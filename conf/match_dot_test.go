package conf

import (
	"testing"
)

func TestPathMatch_dot_Match(t *testing.T) {

	type args struct {
		reg  string
		path string
		spl  string
	}
	type want struct {
		match   bool
		pattern string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "1.1. 常规匹配", args: args{reg: ".a", path: ".a", spl: "."}, want: want{match: true}},
		{name: "1.2. 常规匹配", args: args{reg: ".a.b", path: ".a.b", spl: "."}, want: want{match: true}},
		{name: "1.3. 常规匹配", args: args{reg: ".a.b.c", path: ".a.b.c", spl: "."}, want: want{match: true}},

		{name: "2.1. .**", args: args{reg: ".**", path: ".aa", spl: "."}, want: want{match: true}},

		{name: "2.1. .**", args: args{reg: ".**", path: ".aa", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.bb", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.bb.cc", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.bb.cc.dd", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.js", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.css", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.bb.js", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.bb.css", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.xx.bb.js", spl: "."}, want: want{match: true}},
		{name: "2.1. .**", args: args{reg: ".**", path: ".aa.xx.bb.css", spl: "."}, want: want{match: true}},

		{name: "3.1. .**.a", args: args{reg: ".**.a", path: ".a", spl: "."}, want: want{match: true}},
		{name: "3.1. .**.a", args: args{reg: ".**.a", path: ".aa.bb.a", spl: "."}, want: want{match: true}},
		{name: "3.1. .**.a", args: args{reg: ".**.a", path: ".aa.bb.cc.a", spl: "."}, want: want{match: true}},
		{name: "3.1. .**.a", args: args{reg: ".**.a", path: ".aa.bb.cc.dd.a", spl: "."}, want: want{match: true}},

		{name: "4.1. .**.a.b.js", args: args{reg: ".**.a.b.js", path: ".a.b.js", spl: "."}, want: want{match: true}},
		{name: "4.1. .**.a.b.js", args: args{reg: ".**.a.b.js", path: ".aa.a.b.js", spl: "."}, want: want{match: true}},
		{name: "4.1. .**.a.b.js", args: args{reg: ".**.a.b.js", path: ".aa.bb.a.b.js", spl: "."}, want: want{match: true}},
		{name: "4.1. .**.a.b.js", args: args{reg: ".**.a.b.js", path: ".aa.bb.cc.a.b.js", spl: "."}, want: want{match: true}},

		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".a.b.js", spl: "."}, want: want{match: true}},
		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".aa.a.b.js", spl: "."}, want: want{match: true}},
		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".aa.bb.a.b.js", spl: "."}, want: want{match: true}},
		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".aa.bb.cc.a.b.js", spl: "."}, want: want{match: true}},
		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".aa.bb.cc.a.c.js", spl: "."}, want: want{match: true}},
		{name: "5.1. .**.a.*.js", args: args{reg: ".**.a.*.js", path: ".aa.bb.cc.a.c.css", spl: "."}, want: want{match: false}},

		{name: "6.1. .**.a.b", args: args{reg: ".**.a.b", path: ".a.b", spl: "."}, want: want{match: true}},
		{name: "6.1. .**.a.b", args: args{reg: ".**.a.b", path: ".aa.a.b", spl: "."}, want: want{match: true}},
		{name: "6.1. .**.a.b", args: args{reg: ".**.a.b", path: ".aa.bb.a.b", spl: "."}, want: want{match: true}},
		{name: "6.1. .**.a.b", args: args{reg: ".**.a.b", path: ".aa.bb.cc.a.b", spl: "."}, want: want{match: true}},

		{name: "7.1. .**.a.b.c", args: args{reg: ".**.a.b.c", path: ".a.b.c", spl: "."}, want: want{match: true}},
		{name: "7.1. .**.a.b.c", args: args{reg: ".**.a.b.c", path: ".aa.a.b.c", spl: "."}, want: want{match: true}},
		{name: "7.1. .**.a.b.c", args: args{reg: ".**.a.b.c", path: ".aa.bb.a.b.c", spl: "."}, want: want{match: true}},
		{name: "7.1. .**.a.b.c", args: args{reg: ".**.a.b.c", path: ".aa.bb.cc.a.b.c", spl: "."}, want: want{match: true}},

		{name: "8.1. .*", args: args{reg: ".*", path: ".a", spl: "."}, want: want{match: true}},
		{name: "8.1. .*", args: args{reg: ".*", path: ".aa", spl: "."}, want: want{match: true}},
		{name: "8.1. .*", args: args{reg: ".*", path: ".bb.js", spl: "."}, want: want{match: false}},
		{name: "8.1. .*", args: args{reg: ".*", path: ".cc.css", spl: "."}, want: want{match: false}},
		{name: "8.1. .*", args: args{reg: ".*", path: ".dd.a", spl: "."}, want: want{match: false}},

		{name: "9.1. .*.js", args: args{reg: ".*.js", path: ".a", spl: "."}, want: want{match: false}},
		{name: "9.1. .*.js", args: args{reg: ".*.js", path: ".aa", spl: "."}, want: want{match: false}},
		{name: "9.1. .*.js", args: args{reg: ".*.js", path: ".bb.js", spl: "."}, want: want{match: true}},
		{name: "9.1. .*.js", args: args{reg: ".*.js", path: ".cc.css", spl: "."}, want: want{match: false}},
		{name: "9.1. .*.js", args: args{reg: ".*.js", path: ".dd.a", spl: "."}, want: want{match: false}},

		{name: "10.1. .**.a.**", args: args{reg: ".**.a.**", path: ".a", spl: "."}, want: want{match: true}},
		{name: "10.1. .**.a.**", args: args{reg: ".**.a.**", path: ".aa.a", spl: "."}, want: want{match: true}},
		{name: "10.1. .**.a.**", args: args{reg: ".**.a.**", path: ".aa.bb.a", spl: "."}, want: want{match: true}},
		{name: "10.1. .**.a.**", args: args{reg: ".**.a.**", path: ".a.bb", spl: "."}, want: want{match: true}},
		{name: "10.1. .**.a.**", args: args{reg: ".**.a.**", path: ".xx.yy.a.zz.qq", spl: "."}, want: want{match: true}},

		{name: "11.1. .**.a.*.**", args: args{reg: ".**.a.*.**", path: ".a.b", spl: "."}, want: want{match: true}},
		{name: "11.1. .**.a.*.**", args: args{reg: ".**.a.*.**", path: ".aa.a.b", spl: "."}, want: want{match: true}},
		{name: "11.1. .**.a.*.**", args: args{reg: ".**.a.*.**", path: ".aa.bb.a.b", spl: "."}, want: want{match: true}},
		{name: "11.1. .**.a.*.**", args: args{reg: ".**.a.*.**", path: ".a.b.bb", spl: "."}, want: want{match: true}},
		{name: "11.1. .**.a.*.**", args: args{reg: ".**.a.*.**", path: ".xx.yy.a.b.zz.qq", spl: "."}, want: want{match: true}},

		{name: "12.1. .**.*.a.**.a", args: args{reg: ".**.*.a.**.a", path: ".b.a.a", spl: "."}, want: want{match: true}},
		{name: "12.1. .**.*.a.**.a", args: args{reg: ".**.*.a.**.a", path: ".aa.b.a.a", spl: "."}, want: want{match: true}},
		{name: "12.1. .**.*.a.**.a", args: args{reg: ".**.*.a.**.a", path: ".aa.bb.b.a.xx.a", spl: "."}, want: want{match: true}},
		{name: "12.1. .**.*.a.**.a", args: args{reg: ".**.*.a.**.a", path: ".b.a.bb.a", spl: "."}, want: want{match: true}},
		{name: "12.1. .**.*.a.**.a", args: args{reg: ".**.*.a.**.a", path: ".xx.yy.b.a.zz.qq.a", spl: "."}, want: want{match: true}},

		{name: "13.1. .**.*.a.**.a.*", args: args{reg: ".**.*.a.**.a.*", path: ".b.a.a.c", spl: "."}, want: want{match: true}},
		{name: "13.1. .**.*.a.**.a.*", args: args{reg: ".**.*.a.**.a.*", path: ".aa.b.a.a.c", spl: "."}, want: want{match: true}},
		{name: "13.1. .**.*.a.**.a.*", args: args{reg: ".**.*.a.**.a.*", path: ".aa.bb.b.a.a.c", spl: "."}, want: want{match: true}},
		{name: "13.1. .**.*.a.**.a.*", args: args{reg: ".**.*.a.**.a.*", path: ".b.a.a.c", spl: "."}, want: want{match: true}},
		{name: "13.1. .**.*.a.**.a.*", args: args{reg: ".**.*.a.**.a.*", path: ".xx.yy.b.a.zz.qq.a.c", spl: "."}, want: want{match: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewPathMatch(tt.args.reg)
			match, pattern := m.Match(tt.args.path, tt.args.spl)
			if tt.want.match != match {
				t.Errorf("name:%s,expectMatch:%v,actual:%v,pattern:%s,path:%s", tt.name, tt.want.match, match, pattern, tt.args.path)
			}
		})
	}
}
