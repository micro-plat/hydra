package global

import (
	"testing"
)

func Test_global_Close(t *testing.T) {
	type fields struct {
		isClose bool
		close   chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1-重复调用",
			fields: fields{
				isClose: true,
				close:   make(chan struct{}),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &global{
				isClose: tt.fields.isClose,
				close:   tt.fields.close,
			}
			m.Close()
			m.Close()
		})
	}
}
