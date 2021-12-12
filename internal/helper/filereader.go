package helper

import (
	"go-deduplication/internal/systems"
	"io"
)

type Processor interface {
	Process(fileName string, positions [][]int64, size int64) ([]byte, error)
}

type FileProcessor struct {
	fileReader systems.Reader
}

func NewFileProcessor(fileReader systems.Reader) *FileProcessor {
	return &FileProcessor{fileReader: fileReader}
}

func (p *FileProcessor) Process(fileName string, positions [][]int64, size int64) ([]byte, error) {
	var dataLen int
	for _, position := range positions {
		dataLen += len(position)
	}
	file, err := p.fileReader.OpenFile(fileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	data := make([]byte, int64(dataLen)*size, int64(dataLen)*size)
	var fullTextSize int
	for i := 0; i < len(positions); i++ {
		readData, err := p.fileReader.Read(size, file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		pos := positions[i]
		for _, p := range pos {
			k := 0
			for j := p * size; j < p*size+size; j++ {
				if readData[k] == 0 {
					break
				}
				fullTextSize++
				data[j] = readData[k]
				k++
			}
		}
	}
	result := make([]byte, fullTextSize, fullTextSize)
	for i := 0; i < fullTextSize; i++ {
		result[i] = data[i]
	}
	return result, nil
}
