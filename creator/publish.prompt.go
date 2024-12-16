package creator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/micro-plat/lib4go/types"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type item struct {
	Name   string
	Detail string
	Path   string
}

func newPrompt(path, v string, s interface{}) promptui.Select {
	items := []item{
		{Name: "不更新", Detail: getDetail(v, s), Path: path},
		{Name: "更新", Detail: getDetail(v, s), Path: path},
		{Name: "更新所有配置(存在则更新，不存在则添加)", Path: "所有节点"},
		{Name: "保留已有配置，不进行任何操作", Path: "所有节点"},
	}
	templates := &promptui.SelectTemplates{
		Label:    "节点 {{ .|cyan }} 配置已存在，是否更新？",
		Active:   "\U0001F336 {{ .Name | red }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | red | cyan }} ({{ .Path | red }})",
		Details:  "{{if .Detail}}\n---------- 节点配置参数 ----------{{ .Detail }}{{end}}",
	}
	searcher := func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	return promptui.Select{
		Label:     path,
		Items:     items,
		Templates: templates,
		HideHelp:  true,
		Searcher:  searcher,
	}
}

func getDetail(source string, value interface{}) string {
	target := ""
	if _, ok := value.(string); ok {
		target = value.(string)
	} else {
		s, _ := json.Marshal(value)
		target = string(s)
	}
	s1, s2 := diff(source, target)
	width, _ := terminalWidth()
	wrapper := wrapper(width, true)
	s1 = wrapper(s1)
	s2 = wrapper(s2)
	return fmt.Sprintf("\n存在配置:\n%s\n更新配置:\n%s\n", s1, s2)

}

func diff(source string, target string) (string, string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(source, target, false)
	var buff1 bytes.Buffer
	var buff2 bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			_, _ = buff1.WriteString("\x1b[31m")
			_, _ = buff1.WriteString(text)
			_, _ = buff1.WriteString("\x1b[0m")
		case diffmatchpatch.DiffInsert:
			_, _ = buff2.WriteString("\x1b[32m")
			_, _ = buff2.WriteString(text)
			_, _ = buff2.WriteString("\x1b[0m")
		case diffmatchpatch.DiffEqual:
			_, _ = buff1.WriteString(text)
			_, _ = buff2.WriteString(text)
		}
	}
	return buff1.String(), buff2.String()
}

func terminalWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	s := strings.Split(strings.TrimSuffix(string(out), "\n"), " ")
	return types.GetInt(s[1]), nil
}

type wrapperFunc func(string) string

func wrapper(limit int, breakWords bool) wrapperFunc {
	if limit < 1 {
		panic("Wrapper limit cannot be less than 1.")
	}

	return func(input string) string {
		var wrapped string

		// Split string into array of words
		words := strings.Fields(input)

		if len(words) == 0 {
			return wrapped
		}

		remaining := limit

		if breakWords {
			words = doBreakWords(words, limit)
		}

		for _, word := range words {
			if len(word)+1 > remaining {
				if len(wrapped) > 0 {
					wrapped += "\n"
				}

				wrapped += word
				remaining = limit - len(word)
			} else {
				if len(wrapped) > 0 {
					wrapped += " "
				}

				wrapped += word
				remaining = remaining - (len(word) + 1)
			}
		}

		return wrapped
	}
}

// Break up any words in a given array of words that exceed the given limit.
func doBreakWords(words []string, limit int) []string {
	var result []string

	for _, word := range words {
		if len(word) > limit {
			var parts []string
			var partBuf bytes.Buffer

			for _, char := range word {
				atLimit := partBuf.Len() == limit

				if atLimit {
					parts = append(parts, partBuf.String())

					partBuf.Reset()
				}

				partBuf.WriteRune(char)
			}

			if partBuf.Len() > 0 {
				parts = append(parts, partBuf.String())
			}

			for _, part := range parts {
				result = append(result, part)
			}
		} else {
			result = append(result, word)
		}
	}

	return result
}
