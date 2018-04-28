package mqc

import (
	"encoding/json"
	"fmt"
	"strings"
)

type mqRequest struct {
	service string
	method  string
	raw     string
	form    map[string]interface{}
	header  map[string]string
}

func newMQRequest(service, method, raw string) *mqRequest {
	r := &mqRequest{
		service: service,
		method:  method,
		header:  make(map[string]string),
		form:    make(map[string]interface{}),
		raw:     raw,
	}
	json.Unmarshal([]byte(r.raw), &r.form)
	r.form["__body_"] = r.raw
	return r
}

func (m *mqRequest) GetService() string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(m.service, "/"))
}
func (m *mqRequest) GetMethod() string {
	return m.method
}
func (m *mqRequest) GetForm() map[string]interface{} {
	return m.form
}
func (m *mqRequest) GetHeader() map[string]string {
	return m.header
}
