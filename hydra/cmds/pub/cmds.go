package pub

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/transform"
)

type cmd string

var cmdCD cmd = "cd @path"
var cmdMkdir cmd = "mkdir @path"
var cmdRunScript cmd = "sh ./@path"

func (c cmd) CMD(path string) string {
	ps := make([]interface{}, 0, 2)
	ps = append(ps, "name")
	ps = append(ps, global.AppName)
	return transform.Translate(string(c), "name", global.AppName, "path", path)
}
