package uuid

import (
	"fmt"
	"strconv"
)

//UUID 全局唯一id
type UUID int64

//ToString 字符串
func (u UUID) ToString(prefix ...interface{}) string {
	return fmt.Sprintf("%s%d", fmt.Sprint(prefix...), u)
}

//To16 转16进制字符串
func (u UUID) To16(prefix ...interface{}) string {
	return fmt.Sprintf("%s%x", fmt.Sprint(prefix...), u)
}

//To36 转36进制字符串
func (u UUID) To36(prefix ...interface{}) string {
	return fmt.Sprintf("%s%s", fmt.Sprint(prefix...), intTo36(int64(u)))
}

var charMap map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}

func intTo36(num int64) string {
	n := ""
	var remain int
	var r string
	for num != 0 {
		remain = int(num % 36)
		if 76 > remain && remain > 9 {
			r = charMap[remain]
		} else {
			r = strconv.Itoa(remain)
		}
		n = r + n
		num = num / 36
	}
	return n
}
