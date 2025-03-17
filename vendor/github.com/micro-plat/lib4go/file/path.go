package file

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Exists 检查文件或路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// GetAbs 获取文件绝对路径
func GetAbs(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

// PathMatch 检查路径是否匹配给定的模式（不区分大小写）
func PathMatch(pattern, path string) bool {
	// 将模式转换为正则表达式
	regexPattern := convertPatternToRegex(pattern)
	// 编译正则表达式，并添加 (?i) 标志以支持不区分大小写
	re, err := regexp.Compile("(?i)" + regexPattern)
	if err != nil {
		return false
	}
	// 检查路径是否匹配正则表达式
	return re.MatchString(path)
}

// convertPatternToRegex 将路径模式转换为正则表达式
func convertPatternToRegex(pattern string) string {
	// 转义正则表达式中的特殊字符
	pattern = regexp.QuoteMeta(pattern)
	// 将 ** 替换为匹配任意多段路径的正则表达式
	pattern = strings.ReplaceAll(pattern, "\\*\\*", ".*")
	// 将 * 替换为匹配单段路径的正则表达式
	pattern = strings.ReplaceAll(pattern, "\\*", "[^/]*")
	// 确保路径匹配从开始到结束
	return "^" + pattern + "$"
}
