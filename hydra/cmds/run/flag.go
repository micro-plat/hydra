package run

import (
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/cmds/pkgs"
	"github.com/urfave/cli"
)

//getFlags 获取运行时的参数
func getFlags() []cli.Flag {
	flags := pkgs.GetBaseFlags()
	flags = append(flags, cli.StringFlag{
		Name:        "trace,t",
		Destination: &global.Def.Trace,
		Usage: `-性能跟踪，可选项。用于生成golang的pprof的性能分析数据,支持的模式有:cpu,mem,block,mutex,web。其中web是以http
	 服务的方式提供pprof数据。`,
	})
	flags = append(flags, cli.StringFlag{
		Name:        "traceport,tp",
		Destination: &global.Def.TracePort,
		Usage:       `-性能跟踪，可选项。用于生成golang的pprof的性能分析数据, 当trace为web时候，指定pprof的web端口。默认：19999`,
	})
	flags = append(flags, global.RunCli.GetFlags()...)
	return flags
}
