package conf

import (
	"testing"
)

func TestPathMatch_spec_Match(t *testing.T) {

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
		{name: "1.1 /*", args: args{reg: "/*", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\`, spl: "/"}, want: want{match: true}},
		{name: "1.2 /*/js", args: args{reg: "/*/js", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/js`, spl: "/"}, want: want{match: true}},
		{name: "1.3 /*/js/*", args: args{reg: "/*/js/*", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/js/aaaa`, spl: "/"}, want: want{match: true}},
		{name: "1.4 /*/js/*", args: args{reg: "/*/js/*", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/js/~!@#$%^&*()_+-=<>?:"{}|,.;'[]\`, spl: "/"}, want: want{match: true}},

		{name: "2.1 /**", args: args{reg: "/**", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\`, spl: "/"}, want: want{match: true}},
		{name: "2.2 /**/a", args: args{reg: "/**/a", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/a`, spl: "/"}, want: want{match: true}},
		{name: "2.3 /**/aa/*", args: args{reg: "/**/aa/*", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/aa/bb`, spl: "/"}, want: want{match: true}},
		{name: "2.4 /**/aa/*xyz", args: args{reg: "/**/aa/*xyz", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/aa/bbbbbxyz`, spl: "/"}, want: want{match: true}},
		{name: "2.5 /**/aa/*xyz", args: args{reg: "/**/aa/*xyz", path: `/aa~!@#$%^&*()_+-=<>?:"{}|,.;'[]\/aa/~!@#$%^&*()_+-=<>?:"{}|,.;'[]\xyz`, spl: "/"}, want: want{match: true}},

		{name: `3.1 /*/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, args: args{reg: `/*/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, path: `/aa/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, spl: "/"}, want: want{match: true}},
		{name: `3.2 /**/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, args: args{reg: `/**/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, path: `/aa/~!@#$%^&()_+-=<>?:"{}|,.;'[]\`, spl: "/"}, want: want{match: true}},
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
