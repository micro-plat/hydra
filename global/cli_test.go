package global

import (
	"flag"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func Test_newCli(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *ucli
	}{
		{
			name: "创建cli对象",
			args: args{name: "cli_name"},
			want: &ucli{
				Name:  "cli_name",
				flags: make([]cli.Flag, 0, 1),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCli(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCli() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ucli_AddFlag(t *testing.T) {
	c := newCli("cli_name")
	if c.Name != "cli_name" {
		t.Error("newCli name与预期不匹配", c.Name, "cli_name")
	}

	c.AddFlag("Flag_01", "usage1")

	if len(c.flags) != 1 {
		t.Error("AddFlag 长度与预期不匹配")
	}

	if c.flags[0].GetName() != "Flag_01" {
		t.Error("AddFlag flags.name与预期不匹配", c.flags[0].GetName(), "Flag_01")
	}

	c.AddFlag("Flag_01", "usage1")
	if len(c.flags) != 1 {
		t.Error("AddFlag 添加重复名称，未去重处理", len(c.flags), 1)
	}
}

func Test_ucli_AddSliceFlag(t *testing.T) {
	c := newCli("cli_name")
	if c.Name != "cli_name" {
		t.Error("newCli name与预期不匹配", c.Name, "cli_name")
	}

	c.AddSliceFlag("SliceFlag_01", "usage1")

	if len(c.flags) != 1 {
		t.Error("AddSliceFlag 长度与预期不匹配")
	}

	if c.flags[0].GetName() != "SliceFlag_01" {
		t.Error("AddSliceFlag flags.name与预期不匹配", c.flags[0].GetName(), "SliceFlag_01")
	}

	c.AddSliceFlag("SliceFlag_01", "usage1")
	if len(c.flags) != 1 {
		t.Error("AddSliceFlag 添加重复名称，未去重处理", len(c.flags), 1)
	}
}

func Test_ucli_GetFlags(t *testing.T) {
	type fields struct {
		Name     string
		flags    []cli.Flag
 	}
	tests := []struct {
		name   string
		fields fields
		want   []cli.Flag
	}{
		{
			name: "测试GetFlags-相同的Flag类型",
			fields: fields{
				Name: "getcliflags1",
				flags: []cli.Flag{
					cli.StringFlag{
						Name: "StringFlag_1",
					},
					cli.StringFlag{
						Name: "StringFlag_2",
					},
				},
			},
			want: []cli.Flag{
				cli.StringFlag{
					Name: "StringFlag_1",
				},
				cli.StringFlag{
					Name: "StringFlag_2",
				},
			},
		},
		{
			name: "测试GetFlags-不同的Flag类型",
			fields: fields{
				Name: "getcliflags2",
				flags: []cli.Flag{
					cli.StringFlag{
						Name: "StringFlag_1",
					},
					cli.StringSliceFlag{
						Name: "StringSliceFlag_1",
					},
				},
			},
			want: []cli.Flag{
				cli.StringFlag{
					Name: "StringFlag_1",
				},
				cli.StringSliceFlag{
					Name: "StringSliceFlag_1",
				},
			},
		}, 
 	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ucli{
				Name:     tt.fields.Name,
				flags:    tt.fields.flags,
 			}
			if got := c.GetFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ucli.GetFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doCliCallback(t *testing.T) {
	type args struct {
		c *cli.Context
	}

	app := cli.NewApp()
	app.Name = "Test_AppName"
	flags := []cli.Flag{
		cli.StringFlag{
			Name: "run",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "runtest",
			Usage: "test RUN command",
		},
	}

	app.Flags = flags
	//app.Commands = append(app.Commands, *app.Command("x")

	set := &flag.FlagSet{}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "不存在的command名称-完全包含cmdName",
			args: args{
				c: newCtx(app, set, "runtest"),
			},
			wantErr: true,
		},
		{
			name: "不存在的command名称-只包含cmdName前缀",
			args: args{
				c: newCtx(app, set, "r"),
			},
			wantErr: false,
		},
		{
			name: "不存在的command名称-完全不包含cmdName",
			args: args{
				c: newCtx(app, set, "xxcmd"),
			},
			wantErr: true,
		},
 	}
	t.Log("clis的长度：", len(clis))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("cmd.Name:", tt.args.c.Command.Name)
			if err := doCliCallback(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("doCliCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func newCtx(app *cli.App, set *flag.FlagSet, name string) *cli.Context {
	ctx := cli.NewContext(app, set, nil)
	ctx.Command.Name = name
	return ctx
}
