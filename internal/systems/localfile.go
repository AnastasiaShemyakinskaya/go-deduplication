package systems

import (
	"os"
	"sync"
)

type Writer interface {
	CreateFile(fileName string) (*os.File, error)
	Write(file *os.File, data []byte) error
}

type Reader interface {
	OpenFile(fileName string) (*os.File, error)
	Read(off int64, file *os.File) ([]byte, error)
}

type LocalFile struct {
	mu sync.Mutex
}

func NewLocalFile() *LocalFile {
	return &LocalFile{}
}

func (l *LocalFile) OpenFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (l *LocalFile) CreateFile(fileName string) (*os.File, error) {
	file, err := l.OpenFile(fileName)
	if err != nil {
		file, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func (l *LocalFile) Write(file *os.File, data []byte) error {
	l.mu.Lock()
	_, err := file.Write(data)
	l.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFile) Read(off int64, file *os.File) ([]byte, error) {
	b := make([]byte, off)
	l.mu.Lock()
	_, err := file.Read(b)
	l.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return b, nil
}
