package global

import (
	"flag"
	"testing"

	"github.com/micro-plat/lib4go/assert"

	"github.com/urfave/cli"
)

func Test_newCli(t *testing.T) {
	assert.Equal(t,
		&ucli{
			Name:      "cli_name",
			flags:     make([]cli.Flag, 0, 1),
			flagNames: map[string]bool{},
		},
		newCli("cli_name"),
		"创建cli对象",
	)
}

func Test_ucli_AddFlag(t *testing.T) {
	expectName := "cli_name"
	c := newCli(expectName)

	//测试cli name 与预期是否一致
	assert.Equal(t, expectName, c.Name, "1.1 newCli name与预期不匹配")

	//测试添加flag
	c.AddFlag("-r", "注册中心")
	assert.Equal(t, 1, len(c.flags), "2.1 AddFlag 长度与预期不匹配")
	assert.Equal(t, "-r", c.flags[0].GetName(), "2.2 AddFlag flags.name与预期不匹配")

	//测试重复添加相同的flag
	c.AddFlag("-r", "注册中心")
	assert.Equal(t, 1, len(c.flags), "3.1 AddFlag 添加重复名称，未去重处理")

}

func Test_ucli_AddSliceFlag(t *testing.T) {
	c := newCli("cli_name")

	//测试添加flag
	c.AddSliceFlag("-r", "注册中心")
	assert.Equal(t, 1, len(c.flags), "1.1 AddSliceFlag 长度与预期不匹配")
	assert.Equal(t, "-r", c.flags[0].GetName(), "1.2 AddSliceFlag flags.name与预期不匹配")

	//测试重复添加相同的flag
	c.AddSliceFlag("-r", "注册中心")
	assert.Equal(t, 1, len(c.flags), "2.1 AddSliceFlag 添加重复名称，未去重处理")
}

func Test_ucli_GetFlags(t *testing.T) {
	type fields struct {
		flags []cli.Flag
	}
	tests := []struct {
		name   string
		fields fields
		want   []cli.Flag
	}{
		{name: "1.1 测试GetFlags-相同的Flag类型",
			fields: fields{flags: []cli.Flag{cli.StringFlag{Name: "-r"}, cli.StringFlag{Name: "-c"}}},
			want:   []cli.Flag{cli.StringFlag{Name: "-r"}, cli.StringFlag{Name: "-c"}}},
		{name: "1.2 测试GetFlags-不同的Flag类型",
			fields: fields{flags: []cli.Flag{cli.StringFlag{Name: "-r"}, cli.StringSliceFlag{Name: "-t"}}},
			want:   []cli.Flag{cli.StringFlag{Name: "-r"}, cli.StringSliceFlag{Name: "-t"}}},
	}
	for _, tt := range tests {
		c := &ucli{
			flags: tt.fields.flags,
		}
		assert.Equalf(t, tt.want, c.GetFlags(), tt.name)
	}
}

func Test_doCliCallback(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	//构建 cli 的参数
	app := cli.NewApp()
	app.Name = "Test_AppName"
	flags := []cli.Flag{
		cli.StringFlag{
			Name: "-r",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "test RUN command",
		},
	}

	app.Flags = flags
	set := &flag.FlagSet{}

	tests := []struct {
		name        string
		args        args
		isNillError bool
	}{
		{name: "1. 不存在的command名称-完全包含cmdName", args: args{c: newCtx(app, set, "runtest")}, isNillError: true},
		{name: "2. 不存在的command名称-只包含cmdName前缀", args: args{c: newCtx(app, set, "r")}, isNillError: true},
		{name: "3. 不存在的command名称-完全不包含cmdName", args: args{c: newCtx(app, set, "nonecmd")}, isNillError: true},
	}
	//t.Log("clis的长度：", len(clis))
	for _, tt := range tests {
		//t.Log("cmd.Name:", tt.args.c.Command.Name)
		err := doCliCallback(tt.args.c)
		assert.IsNil(t, tt.isNillError, err, tt.name)
	}
}
func newCtx(app *cli.App, set *flag.FlagSet, name string) *cli.Context {
	ctx := cli.NewContext(app, set, nil)
	ctx.Command.Name = name
	return ctx
}
