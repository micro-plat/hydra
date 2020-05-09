package tpl

import "regexp"

var patterns = map[string]string{
	`where[\t]*[\n]*[\s]*order by`:                   "order by",
	`where[\t]*[\n]*[\s]*group by`:                   "group by",
	`where[\t]*[\n]*[\s]*limit`:                      "limit",
	`where[\t]*[\n]*[\s]*or$`:                        "",
	`where[\t]*[\n]*[\s]*and$`:                       "",
	`where[\t]*[\n]*[\s]*or[\t]*[\n]*[\s]*order by`:  "order by",
	`where[\t]*[\n]*[\s]*or[\t]*[\n]*[\s]*group by`:  "group by",
	`where[\t]*[\n]*[\s]*or[\t]*[\n]*[\s]*limit`:     "limit",
	`where[\t]*[\n]*[\s]*or[\t]*[\n]*[\s]*having`:    "having",
	`where[\t]*[\n]*[\s]*and[\t]*[\n]*[\s]*order by`: "order by",
	`where[\t]*[\n]*[\s]*and[\t]*[\n]*[\s]*group by`: "group by",
	`where[\t]*[\n]*[\s]*and[\t]*[\n]*[\s]*limit`:    "limit",
	`where[\t]*[\n]*[\s]*and[\t]*[\n]*[\s]*having`:   "having",
	`where[\t]*[\n]*[\s]*or[\t|\n|\s]+`:              "where ",
	`where[\t]*[\n]*[\s]*and[\t|\n|\s]+`:             "where ",
	`where[\t]*[\n]*[\s]*$`:                          "",
	`where[\t]*[\n]*[\s]*\)`:                         ")",
}

func replaceSpecialCharacter(s string) string {
	var result = s
	for k, v := range patterns {
		brackets, _ := regexp.Compile(k)
		result = brackets.ReplaceAllStringFunc(result, func(s string) string {
			return v
		})
	}
	return result

}
