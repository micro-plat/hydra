package static

import (
	"strings"
	"time"

	"embed"
)

//DefaultSataticDir 默认静态文件存放路径
const DefaultSataticDir = "./static"

//DefaultHome 默认首页文件名
const DefaultHome = "index.html"

//DefaultExclude 默认需要排除的文件,扩展名,路径
var DefaultExclude = []string{".exe", ".so"}

//DefaultUnrewrite 默认不重写文件
var DefaultUnrewrite = []string{"/favicon.ico", "/robots.txt"}

//Option jwt配置选项
type Option func(*Static)

//WithExclude 排除配置
func WithExclude(excludes ...string) Option {
	return func(s *Static) {
		s.Excludes = excludes
	}
}

//WithUnrewrite 不重写列表
func WithUnrewrite(list ...string) Option {
	return func(s *Static) {
		s.Unrewrites = list
	}
}

//WithAssetsPath 设置资源地址
func WithAssetsPath(path string) Option {
	return func(s *Static) {
		s.Path = path
	}
}

//WithHomePage 设置静首页地址
func WithHomePage(homePage string) Option {
	return func(s *Static) {
		s.HomePage = homePage
	}
}

//WithEmbed 通过嵌入的方式指定压缩文件
func WithEmbed(root string, fs embed.FS) Option {
	return func(s *Static) {
		defEmbedFs[s.serverType] = &embedFs{}
		defEmbedFs[s.serverType].archive = &fs
		defEmbedFs[s.serverType].name = root
	}
}

//WithEmbedBytes 通过嵌入的方式指定压缩文件
func WithEmbedBytes(fileName string, bytes []byte) Option {
	return func(s *Static) {
		defEmbedFs[s.serverType] = &embedFs{}
		defEmbedFs[s.serverType].bytes = bytes
		defEmbedFs[s.serverType].name = fileName
	}
}

//WithAutoRewrite 设置为自动重写
func WithAutoRewrite() Option {
	return func(a *Static) {
		a.AutoRewrite = true
	}
}

//WithDisable 禁用配置
func WithDisable() Option {
	return func(a *Static) {
		a.Disable = true
	}
}

//WithEnable 启用配置
func WithEnable() Option {
	return func(a *Static) {
		a.Disable = false
	}
}

//WithEnableEncryption 启用加密设置
func WithEnableEncryption() Option {
	return func(a *Static) {
		a.EnableEncryption = true
	}
}

//WithEncryptKey 设置 AES-256 加密密钥（推荐32字节）
func WithEncryptKey(key string) Option {
	return func(a *Static) {
		a.EncryptKey = key
	}
}

//WithEncryptIV 设置 AES-GCM 的 IV 偏移量（12字节，前后端约定）
func WithEncryptIV(iv string) Option {
	return func(a *Static) {
		a.EncryptIV = iv
	}
}

//WithEncryptExts 设置需要加密的文件扩展名
func WithEncryptExts(exts ...string) Option {
	return func(a *Static) {
		a.EncryptExts = exts
	}
}

//WithCacheForever 设置永不过期的文件扩展名列表
//d: 过期时间，如 time.Hour*24*30 表示30天
//exts: 文件扩展名，支持两种传参方式：
//   - 多个参数：WithCacheForever(time.Hour*24*30, ".jpg", ".png", ".gif")
//   - 逗号分隔：WithCacheForever(time.Hour*24*30, ".jpg,.png,.gif")
//   - 支持混合使用：WithCacheForever(time.Hour*24*30, ".jpg,.png", ".gif")
func WithCacheForever(d time.Duration, exts ...string) Option {
	return func(s *Static) {
		// 转换时间为秒数
		s.CacheMaxAge = int(d.Seconds())

		// 处理扩展名列表（支持逗号分隔）
		s.CacheForever = parseExtensions(exts)
	}
}

//parseExtensions 解析扩展名列表，支持逗号分隔
func parseExtensions(exts []string) []string {
	var result []string
	for _, ext := range exts {
		// 去除空格
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}
		// 如果包含逗号，则分割
		if strings.Contains(ext, ",") {
			parts := strings.Split(ext, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					result = append(result, part)
				}
			}
		} else {
			result = append(result, ext)
		}
	}
	return result
}
