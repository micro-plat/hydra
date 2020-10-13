package conf

import (
	"testing"
)

func TestMeta_Keys(t *testing.T) {
	tests := []struct {
		name string
		q    Meta
		want map[string]string
	}{
		{
			name: "meta存在单个数据",
			q:    Meta{"key1": "1"},
			want: map[string]string{"key1": "key1"},
		}, {
			name: "meta存在多个数据",
			q:    Meta{"key1": "1", "key2": "2", "key3": "3"},
			want: map[string]string{"key1": "key1", "key2": "key2", "key3": "key3"},
		}, {
			name: "meta存在多个数据,错误返回",
			q:    Meta{"key1": "1", "key2": "2", "key3": "3"},
			want: map[string]string{"key1": "key1", "key3": "key3"},
		}, {
			name: "meta存在多个数据,错误返回1",
			q:    Meta{"key1": "1", "key2": "2", "key3": "3"},
			want: map[string]string{"key1": "key1", "key4": "key4"},
		}, {
			name: "meta不存在数据",
			q:    Meta{},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.q.Keys()
			if len(got) != len(tt.want) {
				t.Errorf("Meta.Keys() = %v, want %v", got, tt.want)
			}
			for _, item := range got {
				if _, ok := tt.want[item]; !ok {
					t.Errorf("Meta.Keys() = %v, want %v", got, tt.want)
				} else {
					delete(tt.want, item)
				}
			}
			if len(tt.want) > 0 {
				t.Errorf("Meta.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_GetString(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		q    Meta
		args args
		want string
	}{
		{
			name: "对象没有数据",
			q:    Meta{},
			args: args{name: "tkey"},
			want: "",
		}, {
			name: "数据不存在",
			q:    Meta{"key1": "1"},
			args: args{name: "tkey"},
			want: "",
		}, {
			name: "数据存在,类型不正确int",
			q:    Meta{"key1": 1},
			args: args{name: "key1"},
			want: "1",
		}, {
			name: "数据存在,类型不正确float",
			q:    Meta{"key1": float32(10.1)},
			args: args{name: "key1"},
			want: "10.1",
		}, {
			name: "数据存在,类型不正确nil",
			q:    Meta{"key1": nil},
			args: args{name: "key1"},
			want: "",
		}, {
			name: "数据存在,类型不正确负数",
			q:    Meta{"key1": -100},
			args: args{name: "key1"},
			want: "-100",
		}, {
			name: "数据存在,类型正确",
			q:    Meta{"key1": "1"},
			args: args{name: "key1"},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.GetString(tt.args.name); got != tt.want {
				t.Errorf("Meta.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_GetInt(t *testing.T) {
	type args struct {
		name string
		def  []int
	}
	tests := []struct {
		name string
		q    Meta
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.GetInt(tt.args.name, tt.args.def...); got != tt.want {
				t.Errorf("Meta.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
