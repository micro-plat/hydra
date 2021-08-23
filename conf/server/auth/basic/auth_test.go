/*
author:taoshouyin
time:2020-10-16
*/

package basic

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func Test_newAuthorization(t *testing.T) {
	tests := []struct {
		name string
		args []*member
		want []*auth
	}{
		{name: "1. 创建空字符串验证对象", args: []*member{&member{UserName: "", Password: ""}}, want: []*auth{&auth{userName: "", auth: createAuth("", "")}}},
		{name: "2. 创建特殊符号的验证对象", args: []*member{&member{UserName: "!@#$%*&^", Password: ")(*^&&^@#$%"}}, want: []*auth{&auth{userName: "!@#$%*&^", password: ")(*^&&^@#$%", auth: createAuth("!@#$%*&^", ")(*^&&^@#$%")}}},
		{name: "3. 创建正常字符串的验证对象", args: []*member{&member{UserName: "taosyteat", Password: "123456"}}, want: []*auth{&auth{userName: "taosyteat", password: "123456", auth: createAuth("taosyteat", "123456")}}},
		{name: "4. 创建多个帐号的验证对象", args: []*member{&member{UserName: "t1", Password: "123"}, &member{UserName: "t1_34", Password: "$%$#12345"}}, want: []*auth{&auth{userName: "t1", password: "123", auth: createAuth("t1", "123")}, &auth{userName: "t1_34", password: "$%$#12345", auth: createAuth("t1_34", "$%$#12345")}}},
	}
	for _, tt := range tests {
		got := newAuthorization(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
