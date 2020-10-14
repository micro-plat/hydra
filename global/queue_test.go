package global

import "testing"

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
		{
			name: "platNameAsPrefix=true,无分隔符",
			fields: fields{
				platNameAsPrefix: true,
				separate:         "",
			},
			args: args{
				n: "myquery:queue",
			},
			want: "testmyquery:queue",
		},
		{
			name: "platNameAsPrefix=true,有分隔符",
			fields: fields{
				platNameAsPrefix: true,
				separate:         ":",
			},
			args: args{
				n: "myquery:queue",
			},
			want: "test:myquery:queue",
		},
		{
			name: "platNameAsPrefix=false,有分隔符",
			fields: fields{
				platNameAsPrefix: false,
				separate:         ":",
			},
			args: args{
				n: "myquery:queue",
			},
			want: "myquery:queue",
		},
		{
			name: "platNameAsPrefix=false,无分隔符",
			fields: fields{
				platNameAsPrefix: false,
				separate:         "",
			},
			args: args{
				n: "myquery:queue",
			},
			want: "myquery:queue",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &messageQueueConf{
				platNameAsPrefix: tt.fields.platNameAsPrefix,
				separate:         tt.fields.separate,
			}
			if got := m.GetQueueName(tt.args.n); got != tt.want {
				t.Errorf("messageQueueConf.GetQueueName() = %v, want %v", got, tt.want)
			}
		})
	}
}
