package obs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/obs/obs"
)

// type Infs interface {
// 	Start() error
// 	Close() error

// 	Save(string, string, []byte) (string, string, error)
// 	Get(string) ([]byte, string, error)
// 	CreateDir(string, string) error
// 	GetFileList(path string, q string, all bool, index int, count int) interface{}

// 	Exists(string) bool
// 	GetDirList(string, int) interface{}
// 	Rename(string, string, string) error
// 	GetScaleImage(root string, path string, width int, height int, quality int) (buff []byte, ctp string, err error)
// 	Conver2PDF(root string, path string) (buff []byte, ctp string, err error)
// 	Registry(tp string)
// }

type OBS struct {
	ak        string
	sk        string
	endpoint  string
	obsClient *obs.ObsClient
}

func NewOBS(ak, sk, endpoint string) *OBS {
	return &OBS{
		ak:       ak,
		sk:       sk,
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
func (o *OBS) getBucketLocation(bucketName string) (string, error) {
	output, err := o.obsClient.GetBucketLocation(bucketName)
	if err != nil {
		return "", err
	}
	return output.Location, nil
}

func (o *OBS) Save(root string, path string, buff []byte) (string, string, error) {
	input := &obs.PutObjectInput{}
	input.Bucket = root
	input.Key = path
	input.Body = bytes.NewReader(buff)

	_, err := o.obsClient.PutObject(input)
	return "", "", err
}
func (o *OBS) Get(path string) ([]byte, string, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = strings.Split(path, "/")[0]
	input.Key = strings.Join(strings.Split(path, "/")[1:], "/")
	output, err := o.obsClient.GetObject(input)
	if err != nil {
		return nil, "", err
	}
	defer output.Body.Close()
	body, err := ioutil.ReadAll(output.Body)
	return body, "", err
}

func (o *OBS) CreateDir(dir string, path string) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	var input = &obs.PutObjectInput{}
	input.Bucket = dir
	input.Key = path
	_, err := o.obsClient.PutObject(input)
	return err
}

func (o *OBS) GetFileList(path string, q string, all bool, index int, count int) interface{} {
	input := &obs.ListObjectsInput{}
	input.Bucket = path
	input.Prefix = q
	output, err := o.obsClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	return nil
}
