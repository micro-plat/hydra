package global

import (
	"fmt"

	"github.com/urfave/cli"
)

//RunCli 执行运行相关的终端参数管理
var RunCli = newCli("run")

//ConfCli 配置处理相关的终端参数
var ConfCli = newCli("conf")

var clis = make([]*ucli, 0, 1)

type CliFlagObject struct {
	RegistryAddr    string
	Name            string
	PlatName        string
	SysName         string
	ServerTypeNames string
	ClusterName     string
}

var FlagVal = &CliFlagObject{}

//IUCLI 终端命令参数
type IUCLI interface {
	AddFlag(name string, usage string) error
	AddSliceFlag(name string, usage string) error
	OnStarting(callback func(ICli) error)
}

type ucli struct {
	Name      string
	flags     []cli.Flag
	flagNames map[string]bool
	callBack  func(ICli) error
}

func newCli(name string) *ucli {
	return &ucli{
		Name:      name,
		flags:     make([]cli.Flag, 0, 1),
		flagNames: map[string]bool{},
	}
}

func (c *ucli) hasFlag(name string) bool {
	if _, ok := c.flagNames[name]; ok {
		return ok
	}
	return false
}

//AddFlag 添加命令行参数
func (c *ucli) AddFlag(name string, usage string) error {
	if c.hasFlag(name) {
		return fmt.Errorf("flag名称%s已存在", name)
	}

	flag := cli.StringFlag{
		Name:  name,
		Usage: usage,
	}
	c.flags = append(c.flags, flag)
	c.flagNames[name] = true
	return nil
}

//AddSliceFlag 添加命令行参数
func (c *ucli) AddSliceFlag(name string, usage string) error {
	if c.hasFlag(name) {
		return fmt.Errorf("flag名称%s已存在", name)
	}
	flag := cli.StringSliceFlag{
		Name:  name,
		Usage: usage,
	}
	c.flags = append(c.flags, flag)
	c.flagNames[name] = true
	return nil
}

//GetFlags 获取可用的flags参数
func (c *ucli) GetFlags() []cli.Flag {
	return c.flags
}

//OnStarting 当启动时执行
func (c *ucli) OnStarting(callback func(ICli) error) {
	c.callBack = callback
}

//Callback 回调onstarting函数
func (c *ucli) Callback(ctx *cli.Context) error {
	if c.callBack == nil {
		return nil
	}
	return c.callBack(ctx)
}

//ICli cli终端参数
type ICli interface {
	IsSet(string) bool
	String(string) string
	StringSlice(string) []string
	FlagNames() []string
	NArg() int
}

func doCliCallback(c *cli.Context) error {
	name := c.Command.FullName()
	for _, cli := range clis {
		if name == cli.Name {
			return cli.Callback(c)
		}
	}
	return nil
}

func init() {
	clis = append(clis, RunCli)
	clis = append(clis, ConfCli)
}
