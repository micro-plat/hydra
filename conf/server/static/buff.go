package static

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const embedArchiveTag = ":EMBED:"

var embedArchive []byte

//saveArchive 保存归档文件
func saveArchive() (string, error) {
	if len(embedArchive) == 0 {
		return "", nil
	}

	rootPath := filepath.Dir(os.Args[0])
	file, err := ioutil.TempFile(rootPath, TempDirName)
	if err != nil {
		return "", err
	}
	_, err = file.Write(embedArchive)
	if err != nil {
		return "", err
	}
	if err = file.Close(); err != nil {
		return "", err
	}
	return file.Name(), nil

}
func removeArchive(f string) {
	os.RemoveAll(f)
}
