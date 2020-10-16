/*
author:taoshouyin
time:2020-10-16
*/

package basic

import (
	"reflect"
	"testing"
)

func Test_newAuthorization(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want []*auth
	}{
		{name: "创建空字符串验证对象", args: args{m: map[string]string{"": ""}}, want: []*auth{&auth{userName: "", auth: createAuth("", "")}}},
		{name: "创建特殊符号的验证对象", args: args{m: map[string]string{"!@#$%*&^": ")(*^&&^@#$%"}}, want: []*auth{&auth{userName: "!@#$%*&^", auth: createAuth("!@#$%*&^", ")(*^&&^@#$%")}}},
		{name: "创建正常字符串的验证对象", args: args{m: map[string]string{"taosyteat": "123456"}}, want: []*auth{&auth{userName: "taosyteat", auth: createAuth("taosyteat", "123456")}}},
		{name: "创建多个帐号的验证对象", args: args{m: map[string]string{"t1": "123", "t1_34": "$%$#12345"}}, want: []*auth{&auth{userName: "t1", auth: createAuth("t1", "123")}, &auth{userName: "t1_34", auth: createAuth("t1_34", "$%$#12345")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAuthorization(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAuthorization() = %v, want %v", got, tt.want)
			}
		})
	}
}
