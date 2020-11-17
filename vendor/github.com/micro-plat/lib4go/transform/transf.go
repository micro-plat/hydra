package transform

//
// type TF map[string]interface{}

// //Translate 翻译字符串kv为map[string]string，map[string]interface{}或以逗号分隔的健值对
// func Translate(format string, kv ...interface{}) string {
// 	if len(kv) == 0 {
// 		panic(fmt.Sprintf("输入的kv必须为：map[string]string，map[string]interface{}，或健值对"))
// 	}
// 	trf := types.NewXMap()
// 	trf.Append(kv...)

// 	brackets, _ := regexp.Compile(`\{@\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*\}`)
// 	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
// 		return trf.GetString(s[2 : len(s)-1])
// 	})
// 	word, _ := regexp.Compile(`@\w+[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*[\.]?\w*`)
// 	result = word.ReplaceAllStringFunc(result, func(s string) string {
// 		return trf.GetString(s[1:])
// 	})
// 	return result
// }
