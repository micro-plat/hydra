package pub

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/transform"
)

type cmd string

var cmdCD cmd = "cd @path"
var cmdMkdir cmd = "mkdir @path"
var cmdRm cmd = "rm -rf @path"
var cmdRunScript cmd = `sh @path @temp_file @bin_name @project_path "@install_params"`

func (c cmd) CMD(path string) string {
	ps := make([]interface{}, 0, 12)
	ps = append(ps, "name")
	ps = append(ps, global.AppName)
	ps = append(ps, "path")
	ps = append(ps, path)
	ps = append(ps, "temp_file")
	ps = append(ps, client.tmpFile)
	ps = append(ps, "bin_name")
	ps = append(ps, client.localPath)
	ps = append(ps, "project_path")
	ps = append(ps, client.projectPath)
	ps = append(ps, "install_params")
	ps = append(ps, runInstall)

	return transform.Translate(string(c), ps...)
}
