package ctx

import (
	"io"
	"os"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
)

var _ context.IFile = &file{}

//file 处理请求上传的文件
type file struct {
	ctx  context.IInnerContext
	meta conf.IMeta
}

//NewFile
func NewFile(ctx context.IInnerContext, meta conf.IMeta) *file {
	return &file{
		ctx:  ctx,
		meta: meta,
	}
}

//SaveFile 保存上传文件到指定路径
func (r *file) SaveFile(fileKey, dst string) error {
	_, reader, _, err := r.ctx.GetFile(fileKey)
	if err != nil {
		return err
	}
	defer reader.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)
	return err
}

//GetFileSize 获取上传文件大小
func (r *file) GetFileSize(fileKey string) (int64, error) {
	_, _, size, err := r.ctx.GetFile(fileKey)
	if err != nil {
		return 0, err
	}
	return size, nil
}

//GetFileName 获取上传文件名称
func (r *file) GetFileName(fileKey string) (string, error) {
	fileName, _, _, err := r.ctx.GetFile(fileKey)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

//GetFileBody 获取上传文件内容
func (r *file) GetFileBody(fileKey string) (io.ReadCloser, error) {
	_, reader, _, err := r.ctx.GetFile(fileKey)
	return reader, err
}
