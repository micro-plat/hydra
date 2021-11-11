package conf

import (
	"testing"
)

func Test_sortString_Len(t *testing.T) {
	tests := []struct {
		name string
		s    sortString
		want int
	}{
		{name: "1. SortStringLen-数据为空数组", s: []string{}, want: 0},
		{name: "2. SortStringLen-数据为空sortString对象", s: sortString{}, want: 0},
		{name: "3. SortStringLen-单个数据对象", s: sortString{"123"}, want: 1},
		{name: "4. SortStringLen-多个数据对象", s: sortString{"123", "3434"}, want: 2},
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
		{name: "1. SortStringSwap-数据为空数组", s: []string{}, args: args{i: 0, j: 0}},
		{name: "2. SortStringSwap-两个数据数据交换", s: sortString{"123", "234"}, args: args{i: 0, j: 1}},
		{name: "3. SortStringSwap-多个数据数据交换", s: sortString{"123", "3434", "656565", "56565444"}, args: args{i: 0, j: 3}},
		{name: "4. SortStringSwap-多个数据数据交换1", s: sortString{"123", "3434", "656565", "56565444", "12222"}, args: args{i: 1, j: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.s) < tt.args.i || len(tt.s) < tt.args.j {
				t.Errorf("用例[%s]out of range", tt.name)
			} else {
				if len(tt.s) == 0 {
					return
				}
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
		{name: "1.1. SortStringLess-单段路径-带*路径排序1", s: sortString{"/t1", "/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "1.2. SortStringLess-单段路径-带*路径排序2", s: sortString{"/*", "/t1"}, args: args{i: 0, j: 1}, want: false},
		{name: "1.3. SortStringLess-单段路径-无*路径排序", s: sortString{"/t1", "/t2"}, args: args{i: 0, j: 1}, want: true},

		{name: "2.1. SortStringLess-多段路径-末尾都带单*", s: sortString{"/t1/*", "/t2/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "2.2. SortStringLess-多段路径-末尾都带单*1", s: sortString{"/t2/*", "/t1/*"}, args: args{i: 0, j: 1}, want: false},
		{name: "2.3. SortStringLess-多段路径-末尾部分带单*", s: sortString{"/t1/*", "/t1/t2"}, args: args{i: 0, j: 1}, want: false},
		{name: "2.4. SortStringLess-多段路径-末尾部分带单*", s: sortString{"/t1/t2", "/t2/*"}, args: args{i: 0, j: 1}, want: true},
		{name: "2.5. SortStringLess-多段路径-末尾**和*进行排序", s: sortString{"/t1/**", "/t1/*"}, args: args{i: 0, j: 1}, want: false}, //**和*号优先级判断
		{name: "2.6. SortStringLess-多段路径-**和精确路径", s: sortString{"/t1/t2", "**"}, args: args{i: 0, j: 1}, want: true},
		{name: "2.7. SortStringLess-多段路径-*和精确路径", s: sortString{"/t1/t2", "*"}, args: args{i: 0, j: 1}, want: true},
		{name: "2.8. SortStringLess-多段路径-两段和三段排序", s: sortString{"/t1/t2/t3", "/t1/t2"}, args: args{i: 0, j: 1}, want: false}, //j字符串被i包含

		{name: "3.1. SortStringLess-ip匹配比较-精确和带*", s: sortString{"192.168.*.*", "192.168.5.94"}, args: args{i: 0, j: 1}, want: false},  //ip中.的判断
		{name: "3.2. SortStringLess-ip匹配比较-中间*和末尾**", s: sortString{"192.168.*.94", "192.168.**"}, args: args{i: 0, j: 1}, want: true}, //ip中.的判断
		{name: "3.3. SortStringLess-ip匹配比较-末尾*和**", s: sortString{"192.168.*.*", "192.168.**"}, args: args{i: 0, j: 1}, want: true},    //ip中.的判断

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
		all []string
	}
	type args struct {
		service string
		seq     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{name: "1.1. PathMatchMatch-精确路径列表-单段", fields: fields{all: []string{"/test1"}}, args: args{service: "/", seq: "/"}, want: false, want1: ""},
		{name: "1.2. PathMatchMatch-精确路径列表-多段列表1", fields: fields{all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3"}}, args: args{service: "/", seq: "/"}, want: true, want1: "/"},
		{name: "1.3. PathMatchMatch-精确路径列表-多段列表2", fields: fields{all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3"}}, args: args{service: "/t2/t3", seq: "/"}, want: false, want1: ""},
		{name: "1.4. PathMatchMatch-精确路径列表-多段列表3", fields: fields{all: []string{"/", "/t1", "/t1/t2", "/t1/t2/t3"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/t2/t3"},
		{name: "1.5. PathMatchMatch-精确路径列表-多段列表，末尾/", fields: fields{all: []string{"/t1/t2"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/t2"},

		{name: "2.1. PathMatchMatch-带*模糊匹配-单段末尾", fields: fields{all: []string{"/t1/*"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/*"},
		{name: "2.2. PathMatchMatch-带*模糊匹配-多段*", fields: fields{all: []string{"/t1/*", "/t1/*/*"}}, args: args{service: "/t1/t3/dd", seq: "/"}, want: true, want1: "/t1/*/*"},
		{name: "2.3. PathMatchMatch-带*模糊匹配-单段*和多段*混合", fields: fields{all: []string{"/t1/*", "/t1/t2/*"}}, args: args{service: "/t1/t2/dd", seq: "/"}, want: true, want1: "/t1/t2/*"},
		{name: "2.4. PathMatchMatch-带*模糊匹配-单段*和精确", fields: fields{all: []string{"/t1/t2/ss", "/t1/t2/*"}}, args: args{service: "/t1/t2/dd", seq: "/"}, want: true, want1: "/t1/t2/*"},
		{name: "2.5. PathMatchMatch-带*模糊匹配-单段*和精确1", fields: fields{all: []string{"/t1/t2/ss", "/t1/t2/*"}}, args: args{service: "/t1/t2/ss", seq: "/"}, want: true, want1: "/t1/t2/ss"},
		{name: "2.6. PathMatchMatch-带*模糊匹配-中间*", fields: fields{all: []string{"/t1/t2/t3", "/t1/*/t4"}}, args: args{service: "/t1/t2/t4", seq: "/"}, want: true, want1: "/t1/*/t4"},
		{name: "2.7. PathMatchMatch-带*模糊匹配-中间*1", fields: fields{all: []string{"/t1/t2/t3", "/t1/*/t4"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/t2/t3"},

		{name: "3.1. PathMatchMatch-单段**", fields: fields{all: []string{"/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/**"},
		{name: "3.2. PathMatchMatch-单段*", fields: fields{all: []string{"/*"}}, args: args{service: "/t1", seq: "/"}, want: true, want1: "/*"},
		{name: "3.3. PathMatchMatch-单段*，不带/", fields: fields{all: []string{"*"}}, args: args{service: "/t1", seq: "/"}, want: false, want1: ""},

		{name: "4.1. PathMatchMatch-*和**混合路径-多段末尾**", fields: fields{all: []string{"/t1/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/**"},
		{name: "4.2. PathMatchMatch-*和**混合路径-末尾**和中间*", fields: fields{all: []string{"/t1/**", "/t1/*/t4"}}, args: args{service: "/t1/t2/t4", seq: "/"}, want: true, want1: "/t1/*/t4"},
		{name: "4.3. PathMatchMatch-*和**混合路径-精确*和**混合", fields: fields{all: []string{"/t1/t2/t3", "/t1/t2/*", "/t1/*/t2", "/t1/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/t2/t3"},
		{name: "4.4. PathMatchMatch-*和**混合路径-*和**混合", fields: fields{all: []string{"/t1/t2/*", "/t1/*/t2", "/t1/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/t2/*"},
		{name: "4.5. PathMatchMatch-*和**混合路径-*和**混合1", fields: fields{all: []string{"/t1/*/t3", "/t1/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/*/t3"},
		{name: "4.6. PathMatchMatch-*和**混合路径-*位置混合", fields: fields{all: []string{"/t1/*/*", "/t1/*"}}, args: args{service: "/t1/t3/dd", seq: "/"}, want: true, want1: "/t1/*/*"},

		{name: "5.1. PathMatchMatch-模糊尾端带/-单*", fields: fields{all: []string{"/t1/*"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/*"},
		{name: "5.2. PathMatchMatch-模糊尾端带/-单*1", fields: fields{all: []string{"/t1/*"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/*"},
		{name: "5.3. PathMatchMatch-模糊尾端带/-尾端**", fields: fields{all: []string{"/t1/**"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/**"},
		{name: "5.4. PathMatchMatch-模糊尾端带/-尾端**1", fields: fields{all: []string{"/t1/**"}}, args: args{service: "/t1/t2", seq: "/"}, want: true, want1: "/t1/**"},
		{name: "5.5. PathMatchMatch-模糊尾端带/-尾端**2", fields: fields{all: []string{"/t1/**"}}, args: args{service: "/t1/t2/t3", seq: "/"}, want: true, want1: "/t1/**"},

		{name: "6.1. PathMatchMatch-ip匹配-精确ip1", fields: fields{all: []string{"192.168.5.124", "192.168.5.22"}}, args: args{service: "192.168.5.94", seq: "."}, want: false, want1: ""},
		{name: "6.2. PathMatchMatch-ip匹配-精确ip2", fields: fields{all: []string{"192.168.5.124", "192.168.5.94"}}, args: args{service: "192.168.5.94", seq: "."}, want: true, want1: "192.168.5.94"},
		{name: "6.3. PathMatchMatch-ip匹配-尾段*模糊匹配", fields: fields{all: []string{"192.168.5.*"}}, args: args{service: "192.168.5.94", seq: "."}, want: true, want1: "192.168.5.*"},
		{name: "6.4. PathMatchMatch-ip匹配-**模糊匹配", fields: fields{all: []string{"192.168.**"}}, args: args{service: "192.168.5.94", seq: "."}, want: true, want1: "192.168.**"},
		{name: "6.5. PathMatchMatch-ip匹配-中段*模糊匹配", fields: fields{all: []string{"192.168.*.94"}}, args: args{service: "192.168.5.94", seq: "."}, want: true, want1: "192.168.*.94"},

		{name: "7.1. PathMatchMatch-路径带:", fields: fields{all: []string{"/t1/t2/:name"}}, args: args{service: "/t1/t2/ss", seq: "/"}, want: false, want1: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewPathMatch(tt.fields.all...)
			got, got1 := a.Match(tt.args.service, tt.args.seq)
			if got != tt.want {
				t.Errorf("PathMatch.Match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PathMatch.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
