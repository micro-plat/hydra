package creator

import (
	"reflect"
	"testing"
)

func Test_newCustomerBuilder(t *testing.T) {
	type args struct {
		s []interface{}
	}
	tests := []struct {
		name string
		args args
		want CustomerBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCustomerBuilder(tt.args.s...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCustomerBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomerBuilder_Sub(t *testing.T) {
	type args struct {
		name string
		s    []interface{}
	}
	tests := []struct {
		name string
		b    CustomerBuilder
		args args
		want ISUB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Sub(tt.args.name, tt.args.s...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerBuilder.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}
