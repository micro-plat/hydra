package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func ScaleImageByPath(p string, width, height, quality int) ([]byte, string, error) {

	basePath := fmt.Sprintf("tmp/scale/%s%s", strings.Replace(p, "/", "|", -1), filepath.Ext(p))
	if fileExist(basePath) {
		fmt.Println("从缓存中读取文件")
		ctp := fmt.Sprintf("image/%s;cache", strings.Trim(filepath.Ext(p), "."))
		data, err := file2Bytes(basePath)
		return data, ctp, err
	}
	fmt.Println("重新构建文件")
	data, err := file2Bytes(p)
	if err != nil {
		return nil, "", fmt.Errorf("读取文件失败:%w %s", err, p)
	}
	buff, ctp, err := ScaleImage(data, width, height, quality)
	if err != nil {
		return nil, "", err
	}
	f, err := os.OpenFile(basePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return buff, ctp, nil
	}
	defer f.Close()
	_, err = f.Write(buff)
	if err != nil {
		fmt.Println(err)
	}
	return buff, ctp, nil
}

func ScaleImage(data []byte, width, height, quality int) ([]byte, string, error) {
	origin, ftp, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 20
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	buf := bytes.Buffer{}
	contentType := "image/jpeg"
	switch strings.ToLower(ftp) {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
		err = jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: quality})
	case "png":
		contentType = "image/png"
		err = png.Encode(&buf, canvas)
	case "gif":
		contentType = "image/gif"
		err = gif.Encode(&buf, canvas, &gif.Options{})
	default:
		return data, contentType, nil
	}

	if err != nil {
		return nil, "", err
	}
	if buf.Len() > len(data) {
		return data, contentType, nil
	}
	return buf.Bytes(), contentType, nil
}
