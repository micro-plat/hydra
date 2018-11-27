package net

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type val struct {
	k string
	v string
}

type Values struct {
	u url.Values
	s []val
}

func NewValues() *Values {
	return &Values{
		u: make(map[string][]string),
		s: make([]val, 0, 2),
	}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (s *Values) Get(key string) string {
	vs := s.u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (s *Values) Set(key, value string) *Values {
	s.u[key] = []string{value}
	s.s = append(s.s, val{k: key, v: value})
	return s
}

//SetMap 指定map[string]interface{}并添加到values中
func (s *Values) SetMap(m map[string]interface{}) *Values {
	for k, v := range m {
		value := fmt.Sprint(v)
		s.u[k] = []string{value}
		s.s = append(s.s, val{k: k, v: value})
		return s
	}
	return s
}

//SetSMap 指定map[string]string并添加到values中
func (s *Values) SetSMap(m map[string]string) *Values {
	for k, value := range m {
		s.u[k] = []string{value}
		s.s = append(s.s, val{k: k, v: value})
		return s
	}
	return s
}

//Encode 参数编号返回编码后的字符串
func (s *Values) Encode() string {
	return s.u.Encode()
}

//Sort 对参数进行参数
func (s *Values) Sort() *Values {
	sort.Slice(s.s, func(i, j int) bool {
		return s.s[i].k < s.s[j].k
	})
	return s
}

//JoinAll 对参数进行拼接
func (s *Values) JoinAll(a string, b string) string {
	buffer := bytes.NewBufferString("")
	for i, v := range s.s {
		buffer.WriteString(v.k)
		buffer.WriteString(a)
		buffer.WriteString(v.v)
		if i < len(s.s)-1 {
			buffer.WriteString(b)
		}
	}
	return strings.Trim(buffer.String(), b)
}

//Join 只拼接值不为空的参数
func (s *Values) Join(a string, b string) string {
	buffer := bytes.NewBufferString("")
	for i, v := range s.s {
		if v.v != "" {
			buffer.WriteString(v.k)
			buffer.WriteString(a)
			buffer.WriteString(v.v)
			if i < len(s.s)-1 {
				buffer.WriteString(b)
			}
		}
	}
	return strings.Trim(buffer.String(), b)
}
