package transform

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

const (
	SymbolEqual = iota
	SymbolNotEqual
	SymbolMore
	SymbolMoreOrEqual
	SymbolLess
	SymbolLessOrEqual
)

//Expression 表达式
type Expression struct {
	Left   string
	Symbol int
	Right  string
}

//Check 翻译模式字符串，并检查值是否符合条件
func Check(query string, trf *Transform) bool {
	exps := Parse(query)
	for _, v := range exps {
		left := trf.Translate(v.Left)
		right := trf.Translate(v.Right)
		switch v.Symbol {
		case SymbolEqual:
			if left != right {
				return false
			}
		case SymbolNotEqual:
			if left == right {
				return false
			}
		case SymbolLess:
			l := types.GetInt(left, 0)
			r := types.GetInt(right, -1)
			if l >= r {
				return false
			}
		case SymbolLessOrEqual:
			l := types.GetInt(left, 0)
			r := types.GetInt(right, -1)
			if l > r {
				return false
			}
		case SymbolMore:
			l := types.GetInt(left, -1)
			r := types.GetInt(right, 0)
			if l <= r {
				return false
			}
		case SymbolMoreOrEqual:
			l := types.GetInt(left, -1)
			r := types.GetInt(right, 0)
			if l < r {
				return false
			}
		}
	}
	return true
}

//Parse 翻译模式字符串，只支持=,!=,>,<
func Parse(query string) (exp []Expression) {
	exp = make([]Expression, 0, 2)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "!="); i >= 0 {
			key, value = key[:i], key[i+2:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolNotEqual, Right: value})
		} else if i := strings.Index(key, ">="); i >= 0 {
			key, value = key[:i], key[i+2:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolMoreOrEqual, Right: value})
		} else if i := strings.Index(key, "<="); i >= 0 {
			key, value = key[:i], key[i+2:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolLessOrEqual, Right: value})
		} else if i := strings.Index(key, "=="); i >= 0 {
			key, value = key[:i], key[i+2:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolEqual, Right: value})
		} else if i := strings.Index(key, ">"); i >= 0 {
			key, value = key[:i], key[i+1:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolMore, Right: value})
		} else if i := strings.Index(key, "<"); i >= 0 {
			key, value = key[:i], key[i+1:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolLess, Right: value})
		} else if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
			exp = append(exp, Expression{Left: key, Symbol: SymbolEqual, Right: value})
		}
	}
	return exp
}
