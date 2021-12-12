package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-deduplication/internal/entity"
	"go-deduplication/internal/helper"
	"go-deduplication/internal/repository"
	"go-deduplication/internal/systems"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"
)

type Restorer interface {
	RestoreFile(ctx context.Context, original []byte, fileName string, path string, function entity.HashFunction, size int) (int64, error)
}

type FileRestorer struct {
	fileHash      repository.FileHashRepository
	fileProcessor helper.Processor
	fileWriter    systems.Writer
}

func NewFileRestorer(
	fileHash repository.FileHashRepository,
	fileProcessor helper.Processor,
	fileWriter systems.Writer,
) *FileRestorer {
	return &FileRestorer{fileHash: fileHash, fileProcessor: fileProcessor, fileWriter: fileWriter}
}

func (r *FileRestorer) RestoreFile(ctx context.Context, original []byte, fileName string, path string, function entity.HashFunction, size int) (int64, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	fileWithoutPath := filepath.Base(fileName)
	name := helper.PrepareFileName(fileWithoutPath, function, size)
	fileWithData := path + "/" + name
	positions, err := r.fileHash.GetPositions(ctx, fileWithData)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("getting positions for file %s", fileName))
	}
	var wrongBytes int64
	data, err := r.fileProcessor.Process(fileWithData, positions, int64(size))
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("getting data for file %s", fileName))
	}
	file, err := r.fileWriter.CreateFile(fileName + "_restored" + strconv.FormatInt(r1.Int63(), 10))
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("creating restored file %s", fileName))
	}
	err = r.fileWriter.Write(file, data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("writing restored file %s", fileName))
	}

	for i := 0; i < len(data) && i < len(original); i++ {
		if original[i] != data[i] {
			wrongBytes++
		}
	}
	return wrongBytes, nil
}
