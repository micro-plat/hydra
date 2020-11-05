package ctx

import (
	"testing"

	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/logger"
)

func Test_response_File(t *testing.T) {
	type fields struct {
		ctx         context.IInnerContext
		conf        app.IAPPConf
		path        *rpath
		raw         rspns
		final       rspns
		noneedWrite bool
		log         logger.ILogger
		asyncWrite  func() error
		specials    []string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &response{
				ctx:         tt.fields.ctx,
				conf:        tt.fields.conf,
				path:        tt.fields.path,
				raw:         tt.fields.raw,
				final:       tt.fields.final,
				noneedWrite: tt.fields.noneedWrite,
				log:         tt.fields.log,
				asyncWrite:  tt.fields.asyncWrite,
				specials:    tt.fields.specials,
			}
			c.File(tt.args.path)
		})
	}
}
