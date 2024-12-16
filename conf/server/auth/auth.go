package auth

var excludes = make([]string, 0, 1)

//AppendExcludes 添加excludes的path
func AppendExcludes(path ...string) {
	excludes = append(excludes, path...)
}

//GetExcludes 获得excludes的path信息
func GetExcludes() []string {
	return excludes
}
