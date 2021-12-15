package internal

import (
	"os"
)

func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	os.Remove(filename)
	return false
}
