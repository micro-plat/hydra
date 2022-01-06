package obs

import (
	"path/filepath"
	"testing"

	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/lib4go/assert"
)

var ak = "HNO8VUJFPF8KLMSHPPQF"
var sk = "RKcWprE1gNORcAukKnhrsnzlPuBZQIdWrb8KL67N"
var bucket = "cdocs-files"
var endpoint = "obs.cn-southwest-2.myhuaweicloud.com"

func TestStart(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)
}
func TestFiles(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)

	input := []struct {
		name string
		buff []byte
	}{
		{name: "test.txt", buff: []byte("test")},
		{name: "a/b/test123.txt", buff: []byte("test123")},
	}
	for _, v := range input {
		p, err := ctx.Save(v.name, v.buff)
		assert.Equal(t, nil, err)
		assert.Equal(t, v.name, p)

		buff, p, err := ctx.Get(v.name)
		assert.Equal(t, nil, err)
		assert.Equal(t, v.buff, buff)

		ok := ctx.Exists(v.name)
		assert.Equal(t, true, ok)

		list := ctx.GetFileList(v.name, infs.GetFullFileName(v.name), true, 0, 10)
		assert.Equal(t, 1, len(list))
		assert.Equal(t, v.name, list[0].Path)

		ctx.Delete(v.name)
		assert.Equal(t, nil, err)

	}
}

func TestDir(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)

	input := []struct {
		name  string
		rname string
		fname string
	}{
		{name: "文件/方案/合同/", rname: "文件/方案/我的合同/", fname: "a.txt"},
	}
	for _, v := range input {
		err = ctx.CreateDir(v.name)
		assert.Equal(t, nil, err)

		list := ctx.GetDirList("/", 2)
		assert.Equal(t, 1, len(list))

		flist := ctx.GetFileList(v.name, "", true, 0, 1)
		assert.Equal(t, 0, len(flist))

		err = ctx.Rename(v.name, v.rname)
		assert.Equal(t, nil, err)

		_, err = ctx.Save(filepath.Join(v.rname, v.fname), []byte("abc222"))
		assert.Equal(t, nil, err)

		flist = ctx.GetFileList(v.rname, "", true, 0, 1)
		assert.Equal(t, 1, len(flist))

		err = ctx.Delete(filepath.Join(v.rname, v.fname))
		assert.Equal(t, nil, err)

		err = ctx.Delete(v.rname)
		assert.Equal(t, nil, err)
	}

}

func TestScaleImage(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)
	buff, ctp, err := ctx.GetScaleImage("照片/我的照片/展展/微信图片_20210514165133.jpg", 80, 80, 80)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", ctp)
	assert.NotEqual(t, 0, len(buff))
}
func TestConver2PDF(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)
	buff, ctp, err := ctx.Conver2PDF("文件/公司文件/合同管理/千行惠捷&聚合动力-渠道推广合作服务协议-2021.07.13.docx")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", ctp)
	assert.NotEqual(t, 0, len(buff))
}
func TestGetFileList(t *testing.T) {
	var ctx = NewOBS(ak, sk, bucket, endpoint, nil)
	err := ctx.Start()
	assert.Equal(t, nil, err)
	list := ctx.GetFileList("文件/公司文件/营销方案", "", true, 0, 100)
	assert.Equal(t, 1, len(list))
}
