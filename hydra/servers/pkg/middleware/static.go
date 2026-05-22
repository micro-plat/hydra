package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/micro-plat/hydra/conf/server/static"
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
			if err := downloadFromHTTP(ctx, static.Path, rpath, static); err != nil {
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

		//设置缓存头
		setCacheHeaders(ctx, static, rpath)

		// 本地文件加密处理
		if static.NeedEncrypt(p) {
			f, ferr := fs.Open(p)
			if ferr != nil {
				ctx.Response().Abort(http.StatusInternalServerError, fmt.Errorf("文件打开失败:%s, 错误:%v", rpath, ferr))
				return
			}
			defer f.Close()

			fileData, rdErr := io.ReadAll(f)
			if rdErr != nil {
				ctx.Response().Abort(http.StatusInternalServerError, fmt.Errorf("文件读取失败:%s, 错误:%v", rpath, rdErr))
				return
			}

			encrypted, encErr := static.DoEncrypt(fileData)
			if encErr != nil {
				ctx.Response().Abort(http.StatusInternalServerError, fmt.Errorf("加密失败:%s, 错误:%v", rpath, encErr))
				return
			}

			ctx.Response().Header("X-Content-Crypto", "aes-256-gcm")
			ctx.Response().Header("X-Content-Crypto-Key", static.KeyFingerprint())
			ctx.Response().Header("Content-Type", "application/octet-stream")
			ctx.Response().Data(http.StatusOK, "application/octet-stream", string(encrypted))
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
func downloadFromHTTP(ctx IMiddleContext, basePath string, rpath string, staticConf *static.Static) error {
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

	// 设置响应头并写入文件内容
	ctx.Response().AddSpecial("static-remote")

	// 设置缓存头
	setCacheHeaders(ctx, staticConf, rpath)

	// 获取密钥指纹
	localFinger := staticConf.KeyFingerprint()

	// 场景1：远程已加密 + 密钥匹配 → 直接透传
	if static.IsRemoteEncrypted(resp.Header, localFinger) {
		ctx.Response().Header("X-Content-Crypto", "aes-256-gcm")
		ctx.Response().Header("X-Content-Crypto-Key", localFinger)
		ctx.Response().Header("Content-Type", "application/octet-stream")
		ctx.Response().Data(http.StatusOK, "application/octet-stream", string(buff))
		return nil
	}

	// 场景2：本地需要加密（配置了密钥 且 文件类型匹配）
	if staticConf.NeedEncrypt(rpath) {
		encrypted, encErr := static.EncryptAndCompress(buff, staticConf.EncryptKey, staticConf.EncryptIV)
		if encErr != nil {
			return fmt.Errorf("加密失败: %v", encErr)
		}
		ctx.Response().Header("X-Content-Crypto", "aes-256-gcm")
		ctx.Response().Header("X-Content-Crypto-Key", localFinger)
		ctx.Response().Header("Content-Type", "application/octet-stream")
		ctx.Response().Data(http.StatusOK, "application/octet-stream", string(encrypted))
		return nil
	}

	// 场景3：不需要加密 → 直接透传
	contentType := resp.Header.Get("Content-Type")

	// 注意：Go 的 http.Client 会自动解压远程返回的 gzip 内容
	// 所以 buff 中存储的是解压后的原始数据
	// 本地的 gzip 中间件会决定是否再次压缩后发送给客户端
	// 不需要手动设置 Content-Encoding

	// 使用 Data 方法直接写入二进制数据
	ctx.Response().Data(http.StatusOK, contentType, string(buff))
	return nil
}

//setCacheHeaders 设置缓存头
func setCacheHeaders(ctx IMiddleContext, staticConf *static.Static, rpath string) {
	if len(staticConf.CacheForever) == 0 {
		return
	}

	// 检查是否匹配永不过期配置
	if shouldCacheForever(staticConf, rpath) {
		ctx.Response().Header("Cache-Control", fmt.Sprintf("public, max-age=%d", staticConf.CacheMaxAge))
	}
}

//shouldCacheForever 判断文件是否应该永不过期
func shouldCacheForever(staticConf *static.Static, rpath string) bool {
	for _, ext := range staticConf.CacheForever {
		if ext == ".*" {
			return true // 所有文件都永不过期
		}
		if strings.HasSuffix(rpath, ext) {
			return true
		}
	}
	return false
}
