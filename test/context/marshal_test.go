package context

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/test/assert"
)

func TestUnmarshalXML(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    map[string]string
		wantErr bool
	}{
		{name: "转换正确的xml", str: `<xml><key>value</key></xml>`, want: map[string]string{"key": "value"}, wantErr: false},
		{name: "转换错误的xml", str: `<xml><key>value</ky></xml>`, want: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		got, err := context.UnmarshalXML(tt.str)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
