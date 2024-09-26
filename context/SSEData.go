package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/micro-plat/lib4go/types"
)

type ISSEData interface {
	Push(data interface{})
	Pop() (bool, string)
	LoopWrite(wr http.ResponseWriter)
	Close()
}

type EventData struct {
	Event string `json:"event,omitempty"`
	Data  string `json:"data"`
}
type SSEData struct {
	data chan interface{}
	done bool
}

func IsSSEData(v interface{}) (bool, ISSEData) {
	if v == nil {
		return false, nil
	}
	resp, ok := v.(ISSEData)
	return ok, resp
}

func NewSSEData(cacheNum ...int) *SSEData {
	return &SSEData{data: make(chan interface{}, types.GetIntByIndex(cacheNum, 0, 32))}
}

func (s *SSEData) Push(data interface{}) {
	if !s.done {
		s.data <- data
	}
}
func (s *SSEData) Close() {
	if !s.done {
		s.done = true
		close(s.data)
	}
}
func (s *SSEData) Pop() (bool, string) {
	if s.done {
		return false, ""
	}
	data, ok := <-s.data
	if !ok {
		return false, ""
	}
	vtpKind := getTypeKind(data)
	if vtpKind == reflect.String {
		return true, fmt.Sprintf("%s", data)
	}
	if buff, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		return true, fmt.Sprintf("%s", string(buff))
	}
}
func (s *SSEData) LoopWrite(wr http.ResponseWriter) {
	wr.Header().Add("Content-Type", UTF8EventStream)
	wr.Header().Add("Cache-Control", "no-cache")
	wr.Header().Add("Connection", "keep-alive")

	for {
		ok, content := s.Pop()
		// fmt.Print(content)
		if !ok {
			return
		}
		if content == "" {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		fmt.Fprintf(wr, "%s\n", content)
		if f, ok := wr.(http.Flusher); ok {
			f.Flush()
		}
	}

}
func getTypeKind(c interface{}) reflect.Kind {
	if c == nil {
		return reflect.String
	}
	value := reflect.ValueOf(c)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.Kind()
}
