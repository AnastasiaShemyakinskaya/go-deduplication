package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-deduplication/internal/entity"
	"go-deduplication/internal/helper"
	"go-deduplication/internal/repository"
	"path/filepath"
)

type Reader interface {
	ReadFile(ctx context.Context, path string, fileName string, function entity.HashFunction, size int) ([]byte, error)
}

type FileReader struct {
	fileHash      repository.FileHashRepository
	fileProcessor helper.Processor
}

func NewFileReader(
	fileHash repository.FileHashRepository,
	fileProcessor helper.Processor,
) *FileReader {
	return &FileReader{fileHash: fileHash, fileProcessor: fileProcessor}
}

func (r *FileReader) ReadFile(ctx context.Context, path string, fileName string, function entity.HashFunction, size int) ([]byte, error) {
	baseFileName := filepath.Base(fileName)
	name := helper.PrepareFileName(baseFileName, function, size)
	fullPath := path + "/" + name
	positions, err := r.fileHash.GetPositions(ctx, fullPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getting positions for file %s", fileName))
	}
	return r.fileProcessor.Process(fullPath, positions, int64(size))
}
