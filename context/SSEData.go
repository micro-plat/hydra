package context

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
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
	done chan struct{}
	once sync.Once
}

func IsSSEData(v interface{}) (bool, ISSEData) {
	if v == nil {
		return false, nil
	}
	resp, ok := v.(ISSEData)
	return ok, resp
}

func NewSSEData(cacheNum ...int) *SSEData {
	return &SSEData{
		data: make(chan interface{}, types.GetIntByIndex(cacheNum, 0, 32)),
		done: make(chan struct{}),
	}
}

func (s *SSEData) Push(data interface{}) {
	if s.done == nil {
		s.done = make(chan struct{})
	}
	select {
	case <-s.done:
		return
	case s.data <- data:
	}
}
func (s *SSEData) Close() {
	if s.done == nil {
		s.done = make(chan struct{})
	}
	s.once.Do(func() {
		close(s.done)
	})
}
func (s *SSEData) Pop() (bool, string) {
	select {
	case data := <-s.data:
		return formatSSEData(data)
	default:
	}
	if s.done == nil {
		s.done = make(chan struct{})
	}
	select {
	case data := <-s.data:
		return formatSSEData(data)
	case <-s.done:
		return false, ""
	}
}

func formatSSEData(data interface{}) (bool, string) {
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
	s.LoopWriteWithContext(context.Background(), wr)
}
func (s *SSEData) LoopWriteWithContext(ctx context.Context, wr http.ResponseWriter) {
	wr.Header().Add("Content-Type", UTF8EventStream)
	wr.Header().Add("Cache-Control", "no-cache")
	wr.Header().Add("Connection", "keep-alive")
	if s.done == nil {
		s.done = make(chan struct{})
	}

	for {
		var data interface{}
		select {
		case <-ctx.Done():
			return
		case data = <-s.data:
		case <-s.done:
			select {
			case data = <-s.data:
			default:
				return
			}
		}
		_, content := formatSSEData(data)
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
