package ctx

import (
	"testing"

	"github.com/micro-plat/hydra/context"
)

func Test_body_GetBody(t *testing.T) {
	type fields struct {
		ctx      context.IInnerContext
		encoding string
		rawBody  bodyValue
		fullBody bodyValue
		mapBody  bodyValue
	}
	tests := []struct {
		name    string
		fields  fields
		wantS   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &body{
				ctx:      tt.fields.ctx,
				encoding: tt.fields.encoding,
				rawBody:  tt.fields.rawBody,
				fullBody: tt.fields.fullBody,
				mapBody:  tt.fields.mapBody,
			}
			gotS, err := w.GetBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("body.GetBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("body.GetBody() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
