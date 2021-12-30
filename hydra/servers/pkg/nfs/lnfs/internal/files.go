package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/**
 * 预览文件相关处理
 */

var (
	officeEtx  = []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
	imageEtx   = []string{".jpeg", ".jpg", ".png", ".gif", ".bmp", ".heic", ".tiff"}
	cadEtx     = []string{".dwg", ".dxf"}
	achieveEtx = []string{".tar.gz", ".tar.bzip2", ".tar.xz", ".zip", ".rar", ".tar", ".7z", "br", ".bz2", ".lz4", ".sz", ".xz", ".zstd"}
	txtEtx     = []string{".txt", ".java", ".php", ".py", ".md", ".js", ".css", ".xml", ".log"}
	videoEtx   = []string{".mp4", ".webm", ".ogg"}
)

const (
	pdfTempRoot  = "./tmp/convert"
	imgThumbnail = "./tmp/thumbnail"
)

func getPDFRootPath(dir string) string {
	return filepath.Join(dir, pdfTempRoot)
}
func getThumbnailPath(dir string, path string) string {
	return filepath.Join(dir, imgThumbnail, strings.Replace(path, "/", "|", -1))
}
func getPDFConverPath(dir string, p string) (dwg string, svg string) {
	name := getFileName(p)
	dwg = filepath.Join(getPDFRootPath(dir), fmt.Sprintf("%s.dwg.pdf", name))
	svg = fmt.Sprintf("%s.svg", p)
	return dwg, svg
}
func getPDFPath(dir string, p string) string {
	return getPDFPathByName(dir, getFileName(p))
}
func getPDFPathByName(dir string, name string) string {
	return filepath.Join(getPDFRootPath(dir), fmt.Sprintf("%s.pdf", name))
}

func fileTypeVerify(url string) (string, string, string) {
	fileName := path.Base(url)  //获取文件名带后缀
	filesuffix := path.Ext(url) //文件后缀

	if strings.Contains(url, ".pdf") {
		return "pdf", ".pdf", fileName
	}

	for _, x := range officeEtx {
		if filesuffix == x {
			return "office", x, fileName
		}
	}

	for _, x := range imageEtx {
		if strings.Contains(url, x) {
			return "image", x, fileName
		}
	}

	for _, x := range cadEtx {
		if strings.Contains(url, x) {
			return "cad", x, fileName
		}
	}

	for _, x := range achieveEtx {
		if strings.Contains(url, x) {
			return "achieve", x, fileName
		}
	}

	for _, x := range txtEtx {
		if strings.Contains(url, x) {
			return "txt", x, fileName
		}
	}

	for _, x := range videoEtx {
		if strings.Contains(url, x) {
			return "video", x, fileName
		}
	}

	return "", "", fileName

}

func ReadFile(filename string) ([]byte, error) {

	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// []byte
	data := make([]byte, stats.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func getFileName(p string) string {
	suffix := path.Base(p)                 //获取文件名带后缀
	ext := path.Ext(suffix)                //获取文件后缀
	return strings.TrimSuffix(suffix, ext) //获取文件名
}
func getExtName(p string) string {
	return strings.Trim(filepath.Ext(p), ".")
}

func checkAndCreateDir(path string) error {
	root := filepath.Dir(path)
	ok, err := pathExists(root)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return os.MkdirAll(root, 0775)

}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		return true
	}
	os.Remove(filename)
	return false
}
