package nfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/conf/server/auth"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/hydra/servers/pkg/nfs/infs"
	"github.com/micro-plat/lib4go/errs"
)

// GetDirList 获取本机目录信息
func (c *cnfs) GetDirList(ctx context.IContext) interface{} {
	return c.infs.GetDirList(infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME,
		ctx.Request().GetString(infs.DIRNAME))),
		ctx.Request().GetInt("deep", 1))
}

// Upload 用户上传文件
func (c *cnfs) Upload(ctx context.IContext) interface{} {
	//读取文件
	name := ctx.Request().GetString(infs.FILENAME, "file")
	name, reader, size, err := ctx.Request().GetFile(name)
	if err != nil {
		return err
	}

	//读取内容
	defer reader.Close()
	buff, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	// 判断是否为 HTTP 远程路径
	if strings.HasPrefix(c.c.Local, "http://") || strings.HasPrefix(c.c.Local, "https://") {
		// 调用 HTTP 远程上传
		return c.uploadToHTTP(ctx, name, buff, size)
	}

	// 保存文件
	path := infs.MultiPath(ctx.Request().Path().Params().GetString("path", ctx.Request().GetString("path")))

	npath, err := c.infs.Save(filepath.Join(path, name), buff)
	if err != nil {
		return err
	}

	// 处理返回结果
	xpath := fmt.Sprintf("%s/%s", strings.TrimRight(c.c.Domain, "/"), strings.Trim(npath, "/"))
	xpath = strings.Replace(xpath, "\\", "/", -1)
	ctx.Response().AddSpecial(fmt.Sprintf("nfs|%s|%d", name, size))
	return map[string]interface{}{
		"path": xpath,
	}
}

// Download 用户下载文件
func (c *cnfs) Download(ctx context.IContext) interface{} {

	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件
	path := filepath.Join(dir, name)
	buff, tp, err := c.infs.Get(path)
	if err != nil {
		return err
	}

	//写入文件
	//未设置文件头
	ctx.Response().ContentType(tp)
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

// GetFileList 获取本机的指定文件的指纹信息，仅master提供对外查询功能
func (c *cnfs) GetFileList(ctx context.IContext) interface{} {
	return c.infs.GetFileList(infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME,
		ctx.Request().GetString(infs.DIRNAME))),
		ctx.Request().GetString("kw"),
		ctx.Request().GetBool("all", false),
		ctx.Request().GetInt("pi", 0),
		ctx.Request().GetInt("ps", 100))
}

// CreateDir 创建目录
func (c *cnfs) CreateDir(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	if dir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.DIRNAME)
	}
	return c.infs.CreateDir(dir)
}

// RenameDir 重命名目录
func (c *cnfs) RenameDir(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	ndir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.NDIRNAME, ctx.Request().GetString(infs.NDIRNAME)))
	if dir == "" || ndir == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\",\":%s\"", infs.DIRNAME, infs.NDIRNAME)
	}
	if dir == ndir {
		return nil
	}
	return c.infs.Rename(dir, ndir)
}

// ImgScale 缩略图生成
func (c *cnfs) ImgScale(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME, ctx.Request().GetString(infs.FILENAME)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件
	path := filepath.Join(dir, name)
	width := ctx.Request().GetInt("w")
	height := ctx.Request().GetInt("h")
	quality := ctx.Request().GetInt("q")
	buff, ctp, err := c.infs.GetScaleImage(path, width, height, quality)
	if err == nil {
		ctx.Response().ContentType(ctp)
		ctx.Response().GetHTTPReponse().Write(buff)
		return nil
	}
	return err
}

// View 获取PDF预览文件
func (c *cnfs) GetPDF4Preview(ctx context.IContext) interface{} {
	//检查参数
	dir := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.DIRNAME, ctx.Request().GetString(infs.DIRNAME)))
	name := infs.MultiPath(ctx.Request().Path().Params().GetString(infs.FILENAME, ctx.Request().GetString(infs.FILENAME)))
	if name == "" {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数不能为空,请求路径中应包含参数 \":%s\"", infs.FILENAME)
	}

	//获取文件
	path := filepath.Join(dir, name)

	buff, _, err := c.infs.GetPDF4Preview(path)
	if err != nil {
		return err
	}
	// ctx.Response().ContentType(contentType) 直接输出流，无需设置contentType
	ctx.Response().GetHTTPReponse().Write(buff)
	return nil
}

func init() {
	auth.AppendExcludes(infs.NOTEXCLUDES...)
}

// uploadToHTTP 上传文件到 HTTP 远程服务
func (c *cnfs) uploadToHTTP(ctx context.IContext, name string, buff []byte, size int64) interface{} {
	// 获取 Processor 配置中的 ServicePrefix
	processor, err := c.app.GetProcessorConf()
	if err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "获取Processor配置失败: %v", err)
	}
	servicePrefix := strings.Trim(processor.ServicePrefix, "/")

	// 构建 HTTP 请求 URL: c.c.Local + servicePrefix + c.c.UploadService
	baseURL := strings.TrimRight(c.c.Local, "/")
	uploadPath := strings.Trim(c.c.UploadService, "/")
	url := fmt.Sprintf("%s/%s/%s", baseURL, servicePrefix, uploadPath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "创建表单文件失败: %v", err)
	}

	// 写入文件内容
	if _, err := part.Write(buff); err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "写入文件内容失败: %v", err)
	}

	// 关闭 writer 以完成 multipart 写入
	if err := writer.Close(); err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "关闭 writer 失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "创建HTTP请求失败: %v", err)
	}

	// 设置 Content-Type 为 multipart/form-data，包含 boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求（设置超时）
	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return errs.NewErrorf(resp.StatusCode, "远程上传失败，URL: %s, 状态码: %d", url, resp.StatusCode)
	}

	// 处理返回结果 - 尝试解析远程服务返回的 JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return errs.NewErrorf(http.StatusInternalServerError, "远程服务返回无效 JSON: %v", err)
	}

	ctx.Response().AddSpecial("nfs-remote")
	return result
}
