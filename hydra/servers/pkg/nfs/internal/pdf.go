package internal

import (
	"fmt"
)

func Conver2PDF(local string) (contentType string, buff []byte, err error) {
	fileType, _, _ := fileTypeVerify(local)
	var resultPath string
	switch fileType {
	case "pdf":
		contentType = "application/x-pdf"
		resultPath = local
	case "image":
		contentType = "image/jpeg"
		resultPath = local
	case "cad":
		contentType = "application/x-pdf"
		resultPath = convertFromCADToPDF(local)
	case "office":
		contentType = "application/x-pdf"
		resultPath = convertToPDF(local)
	case "txt", "":
		contentType = "text/plain"
		resultPath = local
	case "video":
		contentType = "video/mp4"
		resultPath = local
	default:
		return "", nil, fmt.Errorf("文件暂时不支持格式:%s", fileType)
	}
	data, err := file2Bytes(resultPath)
	return "", data, err
}
