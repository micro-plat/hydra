package static

//处理内嵌文件，压缩包则解压并保存到本地，否则通过内存映射

import (
	"embed"
)

type embedFs struct {
	name    string
	archive embed.FS
}

var defEmbedFs = &embedFs{}

//check2FS 检查并转换为fs类型
func (e *embedFs) check2FS() (IFS, error) {
	return newEFS(e.name, e.archive), nil

}
