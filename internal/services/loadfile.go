package services

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/pkg/errors"
	"go-deduplication/internal/entity"
	"go-deduplication/internal/helper"
	"go-deduplication/internal/repository"
	"go-deduplication/internal/systems"
	"hash"
	"os"
	"path/filepath"
)

type Loader interface {
	LoadFile(ctx context.Context, path string, file string, hashFunc entity.HashFunction, byteSize int, data []byte) error
}

type FileLoader struct {
	fileHash   repository.FileHashRepository
	file       repository.FileRepository
	hash       repository.HashRepository
	fileWriter systems.Writer
	db         *systems.Postgres
}

func NewFileLoader(
	fileHash repository.FileHashRepository,
	file repository.FileRepository,
	hash repository.HashRepository,
	fileWriter systems.Writer,
	db *systems.Postgres,
) *FileLoader {
	return &FileLoader{
		fileHash:   fileHash,
		file:       file,
		hash:       hash,
		fileWriter: fileWriter,
		db:         db,
	}
}

func (l *FileLoader) LoadFile(ctx context.Context, path string, file string, hashFunc entity.HashFunction, byteSize int, data []byte) error {
	fileWithoutPath := filepath.Base(file)
	name := helper.PrepareFileName(fileWithoutPath, hashFunc, byteSize)
	fileName := path + "/" + name
	loadFile, err := l.fileWriter.CreateFile(fileName)
	defer loadFile.Close()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating loadFile %s", name))
	}
	fileID, err := l.file.InsertFile(ctx, hashFunc, byteSize, fileName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("inserting loadFile %s", name))
	}
	hashFunction := getHashFunction(hashFunc)
	var position int
	for i := 0; ; {
		if i >= len(data) {
			byteSize = len(data) - (i - byteSize)
			if byteSize == 0 {
				break
			}
			err = l.writeFile(ctx, byteSize, data, i-byteSize, hashFunction, loadFile, fileID, position)
			if err != nil {
				return err
			}
			break
		} else if i+byteSize > len(data) {
			byteSize = len(data) - i
			err = l.writeFile(ctx, byteSize, data, i, hashFunction, loadFile, fileID, position)
			if err != nil {
				return err
			}
			break
		} else {
			err = l.writeFile(ctx, byteSize, data, i, hashFunction, loadFile, fileID, position)
			if err != nil {
				return err
			}
			i = i + byteSize
		}
		position++
	}
	return nil
}

func (l *FileLoader) writeFile(
	ctx context.Context,
	byteSize int,
	data []byte,
	i int,
	hashFunction hash.Hash,
	loadFile *os.File,
	fileID int64,
	position int,
) error {
	bytes := make([]byte, byteSize, byteSize)
	for j := 0; j < byteSize; j++ {
		bytes[j] = data[i+j]
	}
	hashString := hashFunction.Sum(bytes)
	exists, err := l.hash.GetHash(ctx, hashString)
	if err != nil {
		return errors.Wrap(err, "is hash exists query")
	}
	if !exists {
		err := l.fileWriter.Write(loadFile, bytes)
		if err != nil {
			return errors.Wrap(err, "inserting hash")
		}
	}
	err = l.performDbOperations(ctx, hashString, fileID, position)
	if err != nil {
		return errors.Wrap(err, "performing transaction for hash")
	}
	return nil
}

func (l *FileLoader) performDbOperations(ctx context.Context, hashString []byte, fileID int64, i int) error {
	tx, err := l.db.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "opening transaction")
	}
	hashID, err := l.hash.InsertHash(ctx, tx, hashString)
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "inserting hash")
	}

	err = l.fileHash.InsertFileHash(ctx, tx, fileID, hashID, int64(i))
	if err != nil {
		tx.Rollback(ctx)
		return errors.Wrap(err, "inserting file and hash")
	}
	return tx.Commit(ctx)
}

func getHashFunction(function entity.HashFunction) hash.Hash {
	switch function {
	case entity.Md5Function:
		return md5.New()
	case entity.Sha128Function:
		return sha1.New()
	case entity.Sha256Function:
		return sha256.New()
	case entity.Sha512Function:
		return sha512.New()
	}
	return nil
}
