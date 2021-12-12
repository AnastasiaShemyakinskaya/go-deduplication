package helper

import (
	"go-deduplication/internal/entity"
	"strconv"
)

func PrepareFileName(url string, function entity.HashFunction, byteSize int) string {
	return url + "_" + string(function) + "_" + strconv.Itoa(byteSize)
}
