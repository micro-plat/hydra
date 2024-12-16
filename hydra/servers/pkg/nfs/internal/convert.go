package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func convertFromCADToPDF(dir string, filePath string) (string, error) {
	dwg, svg := getPDFConverPath(dir, filePath)
	if fileExist(dwg) {
		return dwg, nil
	}
	commandName := "java"
	params := []string{"-jar", "cad/cad-preview-addon-1.0-SNAPSHOT.jar", "-s", filePath, "-t", svg}

	_, err := executeCMD(commandName, params)
	if err == nil {
		if ok, _ := pathExists(svg); ok {
			return convertToPDF(dir, svg)
		}
	}
	return "", err

}

func convertToPDF(dir string, filePath string) (string, error) {
	pdfPath := getPDFPath(dir, filePath)
	if fileExist(pdfPath) {
		return pdfPath, nil
	}

	rootTemp := getPDFRootPath(dir)
	commandName := ""
	var params []string
	if runtime.GOOS == "windows" {
		commandName = "cmd"
		params = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", rootTemp, filePath}
	} else if runtime.GOOS == "linux" {
		commandName = "libreoffice"
		params = []string{"--invisible", "--headless", "--convert-to", "pdf", "--outdir", rootTemp, filePath}
	} else { // https://ask.libreoffice.org/en/question/12084/how-to-convert-documents-to-pdf-on-osx/
		commandName = "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		params = []string{"--headless", "--convert-to", "pdf", "--outdir", rootTemp, filePath}
	}
	_, err := executeCMD(commandName, params)
	if err != nil {
		return "", err
	}

	paths := strings.Split(path.Base(filePath), ".")

	var tmp = ""
	for i, n := 0, len(paths)-1; i < n; i++ {
		tmp += "." + paths[i]
	}
	resultPath := getPDFPathByName(dir, tmp[1:])
	return resultPath, nil
}

func executeCMD(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	if err != nil {
		return "", fmt.Errorf("命令执行失败：%w %s %s", err, commandName, strings.Join(params, " "))
	}
	return string(buf), nil

}
