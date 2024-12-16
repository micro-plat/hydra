package global

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func Test_messageQueueConf_GetQueueName(t *testing.T) {
	Def.PlatName = "test"

	type fields struct {
		platNameAsPrefix bool
		separate         string
	}
	type args struct {
		n string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "1. platNameAsPrefix=true,无分隔符", fields: fields{platNameAsPrefix: true, separate: ""}, args: args{n: "myquery:queue"}, want: "testmyquery:queue"},
		{name: "2. platNameAsPrefix=true,有分隔符", fields: fields{platNameAsPrefix: true, separate: ":"}, args: args{n: "myquery:queue"}, want: "test:myquery:queue"},
		{name: "3. platNameAsPrefix=false,有分隔符", fields: fields{platNameAsPrefix: false, separate: ":"}, args: args{n: "myquery:queue"}, want: "myquery:queue"},
		{name: "4. platNameAsPrefix=false,无分隔符", fields: fields{platNameAsPrefix: false, separate: ""}, args: args{n: "myquery:queue"}, want: "myquery:queue"}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &messageQueueConf{
			platNameAsPrefix: tt.fields.platNameAsPrefix,
			separate:         tt.fields.separate,
		}
		got := m.GetQueueName(tt.args.n)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
