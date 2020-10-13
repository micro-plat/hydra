package conf

import (
	"reflect"
	"testing"
)

func TestMeta_Keys(t *testing.T) {
	tests := []struct {
		name string
		q    Meta
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Meta.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}
