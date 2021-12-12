package helper

import "io/ioutil"

type AllReader interface {
	ReadAll(fileName string) ([]byte, error)
}

type FileAllReader struct{}

func NewFileAllReader() *FileAllReader {
	return &FileAllReader{}
}

func (f *FileAllReader) ReadAll(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
