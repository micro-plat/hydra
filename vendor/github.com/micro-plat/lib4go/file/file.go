package file

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//CreateFile 根据文件路径(相对或绝对路径)创建文件，如果文件所在的文件夹不存在则自动创建
func CreateFile(path string) (f *os.File, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	dir := filepath.Dir(absPath)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return
	}
	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

//CreateTruncFile 根据文件路径(相对或绝对路径)创建文件，如果文件所在的文件夹不存在则自动创建
func CreateTruncFile(path string) (f *os.File, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	dir := filepath.Dir(absPath)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return
	}
	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}

//SaveBase64Img 保存base64图片
func SaveBase64Img(r string, path string) error {
	//data:image/jpeg;base64,/9j/4QAYRXh
	rows := strings.Split(r, ",")
	content := rows[0]
	if len(rows) > 1 {
		content = rows[1]
	}
	ddd, err := base64.StdEncoding.DecodeString(content) //成图片文件并把文件写入到buffer
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, ddd, 0666)
}
