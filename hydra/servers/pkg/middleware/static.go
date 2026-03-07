package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/services"
)

// Static 静态文件处理插件
func Static() Handler {
	return func(ctx IMiddleContext) {
		//查询静态文件中是否存在
		static, err := ctx.APPConf().GetStaticConf()
		if err != nil {
			ctx.Response().Abort(http.StatusNotExtended, err)
			return
		}
		if static.Disable {
			ctx.Next()
			return
		}

		//处理option请求
		var rpath = ctx.Request().Path().GetRequestPath()
		var method = ctx.Request().Path().GetMethod()
		//是option则处理业务逻辑
		if doOption(ctx, static.Has(rpath)) {
			return
		}
		//优先后端服务调用
		var serverType = ctx.APPConf().GetServerConf().GetServerType()
		if services.Def.Has(serverType, ctx.Request().Path().GetService(), method) {
			ctx.Next()
			return
		}

		//检查请求类型是否为允许的类型
		if !static.AllowRequest(method) {
			ctx.Next()
			return
		}

		// 判断是否为 HTTP 远程路径
		if strings.HasPrefix(static.Path, "http://") || strings.HasPrefix(static.Path, "https://") {
			// 从 HTTP 远程服务下载静态文件
			if err := downloadFromHTTP(ctx, static.Path, rpath); err != nil {
				ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("远程下载文件失败:%s, 错误:%v", rpath, err))
				return
			}
			return
		}

		//读取静态文件
		ctx.Response().AddSpecial("static")
		fs, p, err := static.Get(rpath)
		if err != nil || fs == nil {
			ctx.Response().Abort(http.StatusNotFound, fmt.Errorf("文件不存在:%s", rpath))
			return
		}

		//写入到响应流
		if strings.HasSuffix(p, ".gz") {
			ctx.Response().AddSpecial("gz")
			ctx.Response().Header("Content-Encoding", "gzip")
		}
		ctx.Response().File(p, fs)
	}
}

// downloadFromHTTP 从 HTTP 远程服务下载静态文件
func downloadFromHTTP(ctx IMiddleContext, basePath string, rpath string) error {
	// 构建 HTTP 请求 URL
	baseURL := strings.TrimRight(basePath, "/")
	// 确保路径以 / 开头
	if !strings.HasPrefix(rpath, "/") {
		rpath = "/" + rpath
	}
	url := baseURL + rpath

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("远程服务返回状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 获取 Content-Type
	contentType := resp.Header.Get("Content-Type")
	// 设置响应头并写入文件内容
	ctx.Response().AddSpecial("static-remote")

	// 检查是否为 gzip 压缩文件
	if strings.HasSuffix(rpath, ".gz") {
		ctx.Response().AddSpecial("gz")
		ctx.Response().Header("Content-Encoding", "gzip")
	}

	// 使用 Data 方法直接写入二进制数据
	ctx.Response().Data(http.StatusOK, contentType, string(buff))
	return nil
}
