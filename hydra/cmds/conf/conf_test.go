package conf

import (
	"fmt"
	"os"
	"testing"

	mcli "github.com/micro-plat/cli"
	"github.com/micro-plat/hydra/global"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/urfave/cli"
)

func Test_showNow(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name       string
		args       args
		initParams func()
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	global.ServerTypes = []string{"api"}
	global.FlagVal.ServerTypeNames = "api-mqc"
	app := mcli.New(mcli.WithVersion(global.Version), mcli.WithUsage(global.Usage))

	//ctx := cli.NewContext(app, &flag.FlagSet{}, nil)

	// ctx := cli.NewContext(app, set, nil)
	// ctx.Command.Name = name
	os.Args = []string{"", "conf", "show", "-r", "lm://.", "-p", "xxtest", "-s", "apiserver", "-S", "api", "-c", "c"}

	app.Start()

	fmt.Println(global.Current())
	//assert.Expect(t, err, nil)

	//	fmt.Println("s:", ctx.Args())
	//assert.Expect(t, ctx.Args().Get(0), "abcd")
	//assert.Expect(t, ctx.String("lang"), "spanish")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := showNow(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("showNow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
