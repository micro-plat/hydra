package internal

import (
	"fmt"
)

func Conver2PDF(dir string, path string) (contentType string, buff []byte, err error) {
	fileType, _, _ := fileTypeVerify(path)
	var resultPath string
	switch fileType {
	case "pdf":
		contentType = "application/x-pdf"
		resultPath = path
	case "image":
		contentType = "image/jpeg"
		resultPath = path
	case "cad":
		contentType = "application/x-pdf"
		resultPath, err = convertFromCADToPDF(dir, path)
	case "office":
		contentType = "application/x-pdf"
		resultPath, err = convertToPDF(dir, path)
	case "txt", "":
		contentType = "text/plain"
		resultPath = path
	case "video":
		contentType = "video/mp4"
		resultPath = path
	default:
		return "", nil, fmt.Errorf("文件暂时不支持格式:%s", fileType)
	}
	if err != nil {
		return "", nil, err
	}
	data, err := ReadFile(resultPath)
	if err != nil {
		return "", nil, fmt.Errorf("读取文件失败:%w %s", err, resultPath)
	}
	return "", data, nil
}
