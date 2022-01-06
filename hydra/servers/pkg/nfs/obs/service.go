package obs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/internal"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/obs/obs"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/security/md5"
)

// type Infs interface {
// 	Start() error
// 	Close() error

// 	Save(string, string, []byte) (string, string, error)
// 	Get(string) ([]byte, string, error)
// 	CreateDir(string, string) error
// 	GetFileList(path string, q string, all bool, index int, count int) interface{}
// 	Exists(string) bool
// 	Rename(string, string, string) error

// 	GetDirList(string, int) interface{}

// 	GetScaleImage(root string, path string, width int, height int, quality int) (buff []byte, ctp string, err error)
// 	Conver2PDF(root string, path string) (buff []byte, ctp string, err error)
// 	Registry(tp string)
// }

type OBS struct {
	ak        string
	sk        string
	endpoint  string
	bucket    string
	obsClient *obs.ObsClient
	excludes  []string
	includes  []string
}

func NewOBS(ak, sk, bucket, endpoint string, excludes []string, includes ...string) *OBS {
	obs := &OBS{
		ak:       ak,
		sk:       sk,
		bucket:   bucket,
		endpoint: endpoint,
		excludes: excludes,
		includes: includes,
	}
	return obs
}
func (o *OBS) Start() error {
	var err error
	o.obsClient, err = obs.New(o.ak, o.sk, o.endpoint)
	if err != nil {
		return err
	}
	return nil
}
func (o *OBS) Close() error {
	return nil
}

func (o *OBS) Save(path string, buff []byte) (string, error) {
	input := &obs.PutObjectInput{}
	input.Bucket = o.bucket
	input.Key = path
	input.Body = bytes.NewReader(buff)
	_, err := o.obsClient.PutObject(input)
	return path, err
}
func (o *OBS) Get(path string) ([]byte, string, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = o.bucket
	input.Key = path
	output, err := o.obsClient.GetObject(input)
	if err != nil {
		return nil, "", err
	}
	defer output.Body.Close()
	body, err := ioutil.ReadAll(output.Body)
	tp := infs.GetContentType(path)
	return body, tp, err
}

func (o *OBS) CreateDir(path string) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	var input = &obs.PutObjectInput{}
	input.Bucket = o.bucket
	input.Key = path
	_, err := o.obsClient.PutObject(input)
	return err
}

func (o *OBS) GetFileList(path string, q string, all bool, index int, count int) infs.FileList {
	input := &obs.ListObjectsInput{}
	input.Bucket = o.bucket
	input.Prefix = getCurrentDir(path)
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		return make(infs.FileList, 0)
	}
	list := make(infs.FileList, 0, len(output.Contents))
	for _, val := range output.Contents {
		if strings.HasSuffix(val.Key, "/") || o.fileExclude(val.Key) || !strings.HasPrefix(val.Key, input.Prefix) {
			continue
		}
		if !strings.Contains(infs.GetFullFileName(val.Key), q) {
			continue
		}
		list = append(list, &infs.FileInfo{
			Path:    val.Key,
			DPath:   infs.UnMultiPath(val.Key),
			Name:    infs.GetFileName(val.Key),
			Type:    infs.GetFileType(val.Key),
			ModTime: val.LastModified.Format("2006/01/02 15:04:05"),
			Size:    val.Size,
		})
	}
	sort.Sort(list)

	max := index + count
	if index >= list.Len() {
		return make(infs.FileList, 0)
	}
	if max > list.Len() {
		return list[index:]
	}
	return list[index:max]

}
func (o *OBS) Exists(p string) bool {
	input := &obs.ListObjectsInput{}
	input.Bucket = o.bucket
	input.Prefix = getCurrentDir(p)
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		return false
	}
	for _, v := range output.Contents {
		if v.Key == p {
			return true
		}
	}
	return false
}
func (o *OBS) Delete(path string) error {
	if strings.HasSuffix(path, "/") {
		list := o.GetFileList(path, "", true, 0, 1)
		if len(list) > 0 {
			return fmt.Errorf("目录不为空，不允许删除")
		}
	}

	input := &obs.DeleteObjectInput{}
	input.Bucket = o.bucket
	input.Key = path
	_, err := o.obsClient.DeleteObject(input)
	return err

}
func (o *OBS) Rename(path string, new string) error {

	list := o.GetFileList(path, "", true, 0, 1)
	if len(list) > 0 {
		return fmt.Errorf("目录不为空，不允许重命名")
	}

	input := &obs.CopyObjectInput{}
	input.Bucket = o.bucket
	input.Key = new
	input.CopySourceBucket = o.bucket
	input.CopySourceKey = path
	_, err := o.obsClient.CopyObject(input)
	if err != nil {
		return err
	}
	return o.Delete(path)
}
func (o *OBS) GetDirList(path string, deep int) infs.DirList {
	input := &obs.ListObjectsInput{}
	input.Bucket = o.bucket
	input.Prefix = getCurrentDir(path)
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		return make(infs.DirList, 0)
	}
	list := make(infs.DirList, 0, len(output.Contents))
	for _, val := range output.Contents {
		if val.Key == path || o.fileExclude(path) || !strings.HasSuffix(val.Key, "/") || !strings.HasPrefix(val.Key, path) {
			continue
		}
		list = append(list, &infs.DirInfo{
			ID:      md5.Encrypt(val.Key)[8:16],
			Path:    infs.MultiPath(val.Key),
			DPath:   infs.UnMultiPath(val.Key),
			Name:    infs.GetFileName(val.Key),
			ModTime: val.LastModified.Format("2006/01/02 15:04:05"),
			Size:    val.Size,
			PID:     getParent(val.Key),
		})
	}
	sort.Sort(list)
	return list.GetMultiLevel(path)
}
func getCurrentDir(p string) string {
	if len(p) > 1 {
		if strings.HasSuffix(p, "/") {
			return p
		}
		if !strings.Contains(p, ".") {
			return p + "/"
		}
	}

	f := filepath.Dir(p)
	if f == "." || f == "/" || f == "" {
		return ""
	}
	return f + "/"
}
func (o *OBS) GetScaleImage(path string, width int, height int, quality int) (buff []byte, ctp string, err error) {
	dir := getCurrentDir(path)
	name := infs.GetFullFileName(path)
	thumbnail := filepath.Join(dir, "__thumbnail_"+name)
	if strings.HasPrefix(name, "__thumbnail_") {
		thumbnail = path
	}

	if o.Exists(thumbnail) {
		buff, ctp, err = o.Get(thumbnail)
		return buff, ctp, err
	}
	if path == thumbnail || !o.Exists(path) {
		return nil, "", errs.NewError(404, fmt.Errorf("文件不存在:%s", path))
	}
	buff, ctp, err = o.Get(path)
	buff, err = internal.ScaleImage(buff, width, height, quality)
	if err != nil {
		return nil, "", err
	}
	if _, err = o.Save(thumbnail, buff); err != nil {
		return nil, "", nil
	}

	return buff, ctp, nil
}
func (o *OBS) Conver2PDF(path string) (buff []byte, ctp string, err error) {
	dir := getCurrentDir(path)
	oname := infs.GetFullFileName(path)
	name := infs.GetFileName(path) + ".pdf"

	//构建PDF文件目录
	pdf := filepath.Join(dir, "__pdf_"+name)
	if strings.HasPrefix(name, "__pdf_") {
		pdf = path
	}
	if o.Exists(pdf) {
		buff, ctp, err = o.Get(pdf)
		return buff, ctp, err
	}
	if path == pdf || !o.Exists(path) {
		return nil, "", errs.NewError(404, fmt.Errorf("文件不存在:%s", path))
	}

	//获取原始文件
	buff, ctp, err = o.Get(path)
	if err != nil {
		return nil, "", nil
	}

	//构建临时文件
	tempDir, err := ioutil.TempDir("", "hydad-pdf-")
	if err != nil {
		return nil, "", fmt.Errorf("构建临时目录失败:%w", err)
	}

	//构建临时文件
	file, err := ioutil.TempFile(tempDir, "*"+oname)
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(tempDir)
	defer os.Remove(file.Name())
	defer file.Close()

	file.Write(buff)

	//转换为pdf文件
	fmt.Println("file:", tempDir, file.Name())
	buff, ctp, rpath, err := internal.Conver2PDF(tempDir, file.Name())
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(rpath)

	// if _, err = o.Save(pdf, buff); err != nil {
	// 	return nil, "", nil
	// }

	return buff, ctp, nil
}
func (o *OBS) Registry(tp string) {

}
func getParent(p string) string {
	path := strings.Trim(p, "/")
	if path == "" {
		return path
	}
	items := strings.Split(path, "/")
	return strings.Join(items[0:len(items)-1], "/")
}

func (o *OBS) exclude(p string, excludes ...string) bool {
	nexcludes := append(o.excludes, excludes...)
	return infs.Exclude(p, nexcludes, o.includes...)
}
func (o *OBS) fileExclude(p string) bool {
	nexcludes := append(o.excludes, "__pdf_", "__thumbnail_")
	return infs.Exclude(p, nexcludes, o.includes...)
}
