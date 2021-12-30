package internal

import (
	"fmt"
	"path/filepath"
)

func Conver2PDF(dir string, p string) (contentType string, buff []byte, err error) {
	local := filepath.Join(dir, p)
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
		resultPath, err = convertFromCADToPDF(dir, local)
	case "office":
		contentType = "application/x-pdf"
		resultPath, err = convertToPDF(dir, local)
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
	data, err := ReadFile(resultPath)
	if err != nil {
		return "", nil, fmt.Errorf("读取文件失败:%w %s", err, resultPath)
	}
	return "", data, nil
}
