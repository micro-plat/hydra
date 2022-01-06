package internal

import (
	"fmt"
)

func Conver2PDF(dir string, sourcePath string) (buff []byte, contentType string, p string, err error) {
	fileType, _, _ := fileTypeVerify(sourcePath)
	var resultPath string
	switch fileType {
	case "pdf":
		contentType = "application/x-pdf"
		resultPath = sourcePath
	case "image":
		contentType = "image/jpeg"
		resultPath = sourcePath
	case "cad":
		contentType = "application/x-pdf"
		resultPath, err = convertFromCADToPDF(dir, sourcePath)
	case "office":
		contentType = "application/x-pdf"
		resultPath, err = convertToPDF(dir, sourcePath)
	case "txt", "":
		contentType = "text/plain"
		resultPath = sourcePath
	case "video":
		contentType = "video/mp4"
		resultPath = sourcePath
	default:
		return nil, "", "", fmt.Errorf("文件暂时不支持格式:%s", fileType)
	}
	if err != nil {
		return nil, "", "", err
	}
	data, err := ReadFile(resultPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("读取文件失败:%w %s", err, resultPath)
	}
	return data, contentType, resultPath, nil
}
