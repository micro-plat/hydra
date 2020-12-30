package global

import (
	"fmt"

	"github.com/urfave/cli"
)

//FlagOption 配置选项
type FlagOption func(*ucli)

//WithFlag 设置字符串参数
func WithFlag(name string, usage string) FlagOption {
	return func(c *ucli) {
		if c.hasFlag(name) {
			panic(fmt.Errorf("flag名称%s已存在", name))
		}

		flag := cli.StringFlag{
			Name:  name,
			Usage: usage,
		}
		c.flags = append(c.flags, flag)
		c.flagNames[name] = true
	}
}

//WithFlagByDst 设置字符串参数
func WithFlagByDst(name string, dst *string, usage string) FlagOption {
	return func(c *ucli) {
		if c.hasFlag(name) {
			panic(fmt.Errorf("flag名称%s已存在", name))
		}

		flag := cli.StringFlag{
			Name:        name,
			Destination: dst,
			Usage:       usage,
		}
		c.flags = append(c.flags, flag)
		c.flagNames[name] = true
	}
}

//WithBoolFlag 设置bool参数
func WithBoolFlag(name string, usage string) FlagOption {
	return func(c *ucli) {
		if c.hasFlag(name) {
			panic(fmt.Errorf("flag名称%s已存在", name))
		}

		flag := cli.BoolFlag{
			Name: name,

			Usage: usage,
		}
		c.flags = append(c.flags, flag)
		c.flagNames[name] = true
	}
}

//WithBoolFlagByDst 设置bool参数
func WithBoolFlagByDst(name string, dst *bool, usage string) FlagOption {
	return func(c *ucli) {
		if c.hasFlag(name) {
			panic(fmt.Errorf("flag名称%s已存在", name))
		}

		flag := cli.BoolFlag{
			Name:        name,
			Destination: dst,
			Usage:       usage,
		}
		c.flags = append(c.flags, flag)
		c.flagNames[name] = true
	}
}

//WithSliceFlag 设置数组参数
func WithSliceFlag(name string, usage string) FlagOption {
	return func(c *ucli) {
		if c.hasFlag(name) {
			panic(fmt.Errorf("flag名称%s已存在", name))
		}
		flag := cli.StringSliceFlag{
			Name:  name,
			Usage: usage,
		}
		c.flags = append(c.flags, flag)
		c.flagNames[name] = true
	}
}
