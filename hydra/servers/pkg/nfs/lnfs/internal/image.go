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

func ScaleImageByPath(dir string, p string, width, height, quality int) ([]byte, error) {
	basePath := getThumbnailPath(dir, p)
	if fileExist(basePath) {
		data, err := ReadFile(basePath)
		return data, err
	}
	if err := checkAndCreateDir(basePath); err != nil {
		return nil, fmt.Errorf("创建目录失败:%w", err)
	}

	filePath := filepath.Join(dir, p)
	data, err := ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败:%w %s", err, filePath)
	}
	buff, err := ScaleImage(data, width, height, quality)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(basePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("无法创建文件:%s", basePath)
	}
	defer f.Close()
	_, err = f.Write(buff)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败:%w", err)
	}
	return buff, nil
}

func ScaleImage(data []byte, width, height, quality int) ([]byte, error) {
	origin, ftp, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解码失败:%w", err)
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X / 3
		height = origin.Bounds().Max.Y / 3
	}
	if quality == 0 {
		quality = 40
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	buf := bytes.Buffer{}
	switch strings.ToLower(ftp) {
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, canvas)
	case "gif":
		err = gif.Encode(&buf, canvas, &gif.Options{})
	default:
		return data, nil
	}

	if err != nil {
		return nil, fmt.Errorf("编码失败:%w", err)
	}
	if buf.Len() > len(data) {
		return data, nil
	}
	return buf.Bytes(), nil
}

// func GetContentType(path string) string {
// 	contentType := "image/jpeg"
// 	ftp := getExtName(path)
// 	switch strings.ToLower(ftp) {
// 	case "jpg", "jpeg":
// 		contentType = "image/jpeg"

// 	case "png":
// 		contentType = "image/png"
// 	case "gif":
// 		contentType = "image/gif"
// 	}
// 	return contentType
// }
