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
		resultPath, err = convertFromCADToPDF(local)
	case "office":
		contentType = "application/x-pdf"
		resultPath, err = convertToPDF(local)
	case "txt", "":
		contentType = "text/plain"
		resultPath = local
	case "video":
		contentType = "video/mp4"
		resultPath = local
	default:
		return "", nil, fmt.Errorf("文件暂时不支持格式:%s", fileType)
	}
	if err != nil {
		return "", nil, err
	}
	data, err := file2Bytes(resultPath)
	if err != nil {
		return "", nil, fmt.Errorf("读取文件失败:%w %s", err, resultPath)
	}
	return "", data, nil
}
