package conf

import (
	"sort"
	"testing"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func Test_sortString_Len(t *testing.T) {
	tests := []struct {
		name string
		s    sortString
		want int
	}{
		{name: "t1", s: []string{}, want: 0},
		{name: "t2", s: sortString{}, want: 0},
		{name: "t3", s: sortString{"123"}, want: 1},
		{name: "t4", s: sortString{"123", "3434"}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("sortString.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortString_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    sortString
		args args
	}{
		{name: "t1", s: []string{}, args: args{i: 0, j: 1}},
		{name: "t2", s: sortString{"123", "234"}, args: args{i: 0, j: 1}},
		{name: "t3", s: sortString{"123", "3434", "656565", "56565444"}, args: args{i: 0, j: 3}},
		{name: "t4", s: sortString{"123", "3434", "656565", "56565444", "12222"}, args: args{i: 1, j: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.s) < tt.args.i || len(tt.s) < tt.args.j {
				t.Errorf("用例[%s]out of range", tt.name)
			} else {
				oi := tt.s[tt.args.i]
				oj := tt.s[tt.args.j]
				tt.s.Swap(tt.args.i, tt.args.j)
				if tt.s[tt.args.i] != oj || tt.s[tt.args.j] != oi {
					t.Errorf("用例[%s]数据交换失败", tt.name)
				}
			}
		})
	}
}

func Test_sortString_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    sortString
		args args
		want bool
	}{
		{name: "t1", s: sortString{"/t1", "/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "t2", s: sortString{"/*", "/t1"}, args: args{i: 0, j: 1}, want: false},
		{name: "t3", s: sortString{"/t1", "/t2"}, args: args{i: 0, j: 1}, want: true},
		{name: "t4", s: sortString{"/t1/*", "/t2/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "t5", s: sortString{"/t2/*", "/t1/*"}, args: args{i: 0, j: 1}, want: false},
		{name: "t6", s: sortString{"/t1/*", "/t1/t2"}, args: args{i: 0, j: 1}, want: false},
		{name: "t7", s: sortString{"/t1/t2", "/t2/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "t8", s: sortString{"/t1/**", "/t1/*"}, args: args{i: 0, j: 1}, want: false},              //**和*号优先级判断
		{name: "t9", s: sortString{"/t1/t2/t3", "/t1/t2"}, args: args{i: 0, j: 1}, want: true},           //j字符串被i包含
		{name: "t10", s: sortString{"192.168.*.*", "192.168.5.94"}, args: args{i: 0, j: 1}, want: false}, //ip中.的判断
		{name: "t11", s: sortString{"192.168.*.94", "192.168.**"}, args: args{i: 0, j: 1}, want: true},   //ip中.的判断
		{name: "t12", s: sortString{"192.168.*.*", "192.168.**"}, args: args{i: 0, j: 1}, want: true},    //ip中.的判断
		{name: "t13", s: sortString{"/t1/t2", "/t1/t2/t3"}, args: args{i: 0, j: 1}, want: true},          //i的字符串被j包含时,数组超过限制崩溃
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("sortString.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathMatch_Match(t *testing.T) {
	type fields struct {
		cache cmap.ConcurrentMap
		all   []string
	}
	type args struct {
		service string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
		{
			name:   "t1",
			fields: fields{cache: cmap.New(6), all: []string{"/test1"}},
			args:   args{service: "/"},
			want:   false,
			want1:  "",
		}, {
			name:   "t2",
			fields: fields{cache: cmap.New(6), all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3", "/t1/t2/t3/t4", "/t1/t2/t3/t4/t5"}},
			args:   args{service: "/"},
			want:   true,
			want1:  "/",
		}, {
			name:   "t3",
			fields: fields{cache: cmap.New(6), all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3", "/t1/t2/t3/t4", "/t1/t2/t3/t4/t5"}},
			args:   args{service: "/t2/t3"},
			want:   false,
			want1:  "",
		}, {
			name:   "t4",
			fields: fields{cache: cmap.New(6), all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3", "/t1/t2/t3/t4", "/t1/t2/t3/t4/t5"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/t2/t3",
		}, {
			name:   "t5",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*"}},
			args:   args{service: "/t1/t2"},
			want:   true,
			want1:  "/t1/*",
		}, {
			name:   "t6",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*"}},
			args:   args{service: "/t1/t3"},
			want:   true,
			want1:  "/t1/*",
		}, {
			name:   "t7",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*", "/t1/*/*"}},
			args:   args{service: "/t1/t3/dd"},
			want:   true,
			want1:  "/t1/*/*",
		}, {
			name:   "t8",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*", "/t1/t2/*"}},
			args:   args{service: "/t1/t2/dd"},
			want:   true,
			want1:  "/t1/t2/*",
		}, {
			name:   "t9",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/ss", "/t1/t2/*"}},
			args:   args{service: "/t1/t2/dd"},
			want:   true,
			want1:  "/t1/t2/*",
		}, {
			name:   "t10",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/ss", "/t1/t2/*"}},
			args:   args{service: "/t1/t2/ss"},
			want:   true,
			want1:  "/t1/t2/ss",
		}, {
			name:   "t11",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/:name"}},
			args:   args{service: "/t1/t2/ss"},
			want:   false,
			want1:  "",
		}, {
			name:   "t12",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/t3", "/t1/*/t4"}},
			args:   args{service: "/t1/t2/t4"},
			want:   true,
			want1:  "/t1/*/t4",
		}, {
			name:   "t13",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/t3", "/t1/*/t4"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/t2/t3",
		}, {
			name:   "t14",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/**"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/**",
		}, {
			name:   "t15",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/**", "/t1/*/t4"}},
			args:   args{service: "/t1/t2/t4"},
			want:   true,
			want1:  "/t1/*/t4",
		}, {
			name:   "t16",
			fields: fields{cache: cmap.New(6), all: []string{"192.168.5.124", "192.168.5.22"}},
			args:   args{service: "192.168.5.94"},
			want:   false,
			want1:  "",
		}, {
			name:   "t17",
			fields: fields{cache: cmap.New(6), all: []string{"192.168.5.124", "192.168.5.94"}},
			args:   args{service: "192.168.5.94"},
			want:   true,
			want1:  "192.168.5.94",
		}, {
			name:   "t18",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/t3", "/t1/t2/*", "/t1/*/t2", "/t1/**"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/t2/t3",
		}, {
			name:   "t19",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/t2/*", "/t1/*/t2", "/t1/**"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/t2/*",
		}, {
			name:   "t20",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*/t3", "/t1/**"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/t1/*/t3",
		}, {
			name:   "t21",
			fields: fields{cache: cmap.New(6), all: []string{"/**"}},
			args:   args{service: "/t1/t2/t3"},
			want:   true,
			want1:  "/**",
		}, {
			name:   "t22",
			fields: fields{cache: cmap.New(6), all: []string{"/t1/*/*", "/t1/*"}},
			args:   args{service: "/t1/t3/dd"},
			want:   true,
			want1:  "/t1/*/*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PathMatch{
				cache: tt.fields.cache,
				all:   tt.fields.all,
			}
			sort.Sort(sortString(a.all))
			got, got1 := a.Match(tt.args.service)
			if got != tt.want {
				t.Errorf("PathMatch.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PathMatch.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
