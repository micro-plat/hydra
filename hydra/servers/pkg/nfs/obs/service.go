package obs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/internal"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/obs/obs"
	"github.com/micro-plat/lib4go/errs"
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
}

func NewOBS(ak, sk, bucket, endpoint string) *OBS {
	return &OBS{
		ak:       ak,
		sk:       sk,
		bucket:   bucket,
		endpoint: endpoint,
	}
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
	input.Prefix = getPath(path)
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		return make(infs.FileList, 0)
	}
	list := make(infs.FileList, 0, len(output.Contents))
	for _, val := range output.Contents {
		if strings.HasSuffix(val.Key, "/") {
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
	input.Prefix = getPath(p)
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
	input.Prefix = getPath(path)
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		return make(infs.DirList, 0)
	}
	list := make(infs.DirList, 0, len(output.Contents))
	for _, val := range output.Contents {
		if val.Key == path || !strings.HasSuffix(val.Key, "/") {
			continue
		}
		list = append(list, &infs.DirInfo{
			Path:    val.Key,
			DPath:   infs.UnMultiPath(val.Key),
			Name:    infs.GetFileName(val.Key),
			ModTime: val.LastModified.Format("2006/01/02 15:04:05"),
			Size:    val.Size,
		})
	}
	return list
}
func getPath(p string) string {
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		return p
	}
	f := filepath.Dir(p)
	if f == "." || f == "/" || f == "" {
		return ""
	}
	return f + "/"
}
func (o *OBS) GetScaleImage(path string, width int, height int, quality int) (buff []byte, ctp string, err error) {
	dir := getPath(path)
	name := infs.GetFileName(path)
	thumbnail := filepath.Join(dir, "thumbnail_"+name)
	if o.Exists(thumbnail) {
		buff, ctp, err = o.Get(thumbnail)
		return buff, ctp, err
	}
	if !o.Exists(path) {
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
	return nil, "", nil
}
func (o *OBS) Registry(tp string) {

}
