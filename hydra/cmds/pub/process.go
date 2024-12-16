package pub

import (
	"fmt"
	"strings"
)

// Progress 进度
type Progress int

// Show 显示进度
func (x Progress) Show() {
	percent := int(x)
	total := 50 //格子数

	middle := int(percent * total / 100.0)

	arr := make([]string, total)
	for j := 0; j < total; j++ {
		if j < middle-1 {
			arr[j] = "-"
		} else if j == middle-1 {
			arr[j] = ">"
		} else {
			arr[j] = " "
		}
	}
	bar := fmt.Sprintf("编译文件上传进度：[%s]", strings.Join(arr, ""))
	fmt.Printf("\r%s %d%%", bar, percent)
}
