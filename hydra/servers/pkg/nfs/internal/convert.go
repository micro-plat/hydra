package internal

import (
	"path"
	"runtime"
	"strings"
)

func convertFromCADToPDF(filePath string) (string, error) {
	basePath := "tmp/convert/" + getFileNameOnly(filePath) + ".dwg.pdf"
	if fileExist(basePath) {
		return basePath, nil
	}
	commandName := "java"
	params := []string{"-jar", "cad/cad-preview-addon-1.0-SNAPSHOT.jar", "-s", filePath, "-t", filePath + ".svg"}

	_, err := interactiveToexec(commandName, params)
	if err == nil {
		resultPath := filePath + ".svg"
		if ok, _ := pathExists(resultPath); ok {
			return convertToPDF(resultPath)
		}
	}
	return "", err

}

func convertToPDF(filePath string) (string, error) {

	basePath := "tmp/convert/" + getFileNameOnly(filePath) + ".pdf"
	if fileExist(basePath) {
		return basePath, nil
	}
	commandName := ""
	var params []string
	if runtime.GOOS == "windows" {
		commandName = "cmd"
		params = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", "tmp/convert/", filePath}
	} else if runtime.GOOS == "linux" {
		commandName = "libreoffice"
		params = []string{"--invisible", "--headless", "--convert-to", "pdf", "--outdir", "tmp/convert/", filePath}
	} else { // https://ask.libreoffice.org/en/question/12084/how-to-convert-documents-to-pdf-on-osx/
		commandName = "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		params = []string{"--headless", "--convert-to", "pdf", "--outdir", "tmp/convert/", filePath}
	}
	_, err := interactiveToexec(commandName, params)
	if err != nil {
		return "", err
	}

	paths := strings.Split(path.Base(filePath), ".")

	var tmp = ""
	for i, n := 0, len(paths)-1; i < n; i++ {
		tmp += "." + paths[i]
	}
	resultPath := "tmp/convert/" + tmp[1:] + ".pdf"
	return resultPath, nil

}
