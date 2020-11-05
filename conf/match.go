package conf

import (
	"sort"
	"strings"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type sortString []string

func (s sortString) Len() int { return len(s) }

func (s sortString) Swap(i, j int) {
	if i >= len(s) || j >= len(s) {
		return
	}
	s[i], s[j] = s[j], s[i]
}

func (s sortString) Less(i, j int) bool {
	il := len(s[i])
	jl := len(s[j])
	for x := 0; x < jl && x < il; x++ {
		if s[i][x] == s[j][x] {
			continue
		}
		if s[i][x] == []byte("*")[0] {
			return false
		}
		if s[j][x] == []byte("*")[0] {
			return true
		}
		return s[i][x] < s[j][x]
	}
	return s[i] < s[j]
}

//PathMatch 构建模糊匹配缓存查找管理器
type PathMatch struct {
	cache cmap.ConcurrentMap
	all   []string
}

//NewPathMatch 构建模糊匹配缓存查找管理器
func NewPathMatch(all ...string) *PathMatch {
	i := &PathMatch{
		cache: cmap.New(6),
		all:   all,
	}
	sort.Sort(sortString(i.all))
	return i
}

//Match 是否匹配，支付完全匹配，模糊匹配，分段匹配
func (a *PathMatch) Match(service string, seq ...string) (bool, string) {
	if v, ok := a.cache.Get(service); ok {
		return v != "", v.(string)
	}
	nseq := "/"
	if len(seq) > 0 {
		nseq = seq[0]
	}
	sparties := strings.Split(service, nseq)
	//排除指定请求
	for _, u := range a.all {
		//完全匹配
		if strings.EqualFold(u, service) {
			a.cache.SetIfAbsent(service, u)
			return true, u
		}
		//分段模糊
		uparties := strings.Split(u, nseq)
		//取较少的数组长度
		uc := len(uparties)
		sc := len(sparties)
		/*
			路径处理模式：
			1. /a/b/ *
			2. /a/ **
			3. /a/ * /d

			ip处理模式：
			1. 192.168.0.*
			2. 192.168.**
			3. 192.*.0.1
		**/

		//长度不匹配，且未包含**,跳过
		if uc != sc && !strings.HasSuffix(u, "**") {
			continue
		}

		//原段较长，不可能匹配跳过
		if uc > sc {
			continue
		}

		//原段较短，或有**进行分段检查
		isMatch := true
		for i := 0; i < uc; i++ {

			//此段为 **
			if uparties[i] == "**" {
				a.cache.SetIfAbsent(service, u)
				return true, u
			}

			//此段为 *,匹配后续段
			if uparties[i] == "*" {
				for j := i + 1; j < uc; j++ {
					if uparties[j] == "*" {
						continue
					}
					if uparties[j] != sparties[j] {
						isMatch = false
						break
					}
				}
				if !isMatch {
					break
				}
				a.cache.SetIfAbsent(service, u)
				return true, u
			}
			if uparties[i] != sparties[i] {
				break
			}
		}

	}
	a.cache.SetIfAbsent(service, "")
	return false, ""
}
