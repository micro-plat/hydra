package static

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const embedArchiveTag = ":EMBED:"

var embedArchive []byte
var embedExt string

func getEmbedFileName() string {
	return fmt.Sprintf("%s%s", TempArchiveName, embedExt)
}

//saveArchive 保存归档文件
func saveArchive() (string, error) {
	if len(embedArchive) == 0 {
		return "", nil
	}

	rootPath := filepath.Dir(os.Args[0])
	file, err := ioutil.TempFile(rootPath, getEmbedFileName())
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
