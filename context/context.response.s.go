package context

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

func getJSONContent(t int, data interface{}) (int, interface{}, error) {
	if data == nil {
		return t, nil, nil
	}
	if err, ok := data.(error); ok {
		data = map[string]interface{}{"err": err.Error()}
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.String:
		value := []byte(data.(string))
		switch {
		case (t == CT_JSON || t == CT_DEF) && json.Valid(value) && (bytes.HasPrefix(value, []byte("{")) ||
			bytes.HasPrefix(value, []byte("["))):
			return CT_JSON, data, nil
		case (t == CT_XML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<?xml")):
			return CT_XML, data, nil
		case (t == CT_HTML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<!DOCTYPE html")):
			return CT_HTML, data, nil
		}
		switch {
		case t == CT_JSON || t == CT_DEF:
			return CT_JSON, map[string]interface{}{"data": data}, nil
		default:
			return t, data, nil
		}
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Array:
		switch {
		case t == CT_JSON || t == CT_DEF:
			buff, err := json.Marshal(data)
			if err != nil {
				return t, nil, err
			}
			return CT_JSON, string(buff), nil
		case t == CT_XML:
			buff, err := xml.Marshal(data)
			if err != nil {
				return t, nil, err
			}
			return CT_XML, string(buff), nil
		case t == CT_YMAL:
			buff, err := yaml.Marshal(data)
			if err != nil {
				return CT_YMAL, nil, err
			}
			return CT_YMAL, string(buff), nil
		default:
			return t, fmt.Sprintf("%+v", data), nil
		}

	default:
		switch {
		case t == CT_JSON || t == CT_DEF:
			return getJSONContent(t, map[string]interface{}{"data": data})
		case t == CT_YMAL:
			buff, err := yaml.Marshal(map[string]interface{}{
				"data": data,
			})
			if err != nil {
				return t, nil, err
			}
			return CT_YMAL, string(buff), nil
		default:
			return t, fmt.Sprint(data), nil
		}
	}
}

func (r *Response) GetJSONRenderContent() (int, interface{}, error) {
	return getJSONContent(r.getContentType(), r.GetContent())

}

func (r *Response) GetHTMLRenderContent() (int, interface{}, error) {
	data := r.GetContent()
	t := r.getContentType()
	if data == nil {
		return t, nil, nil
	}
	if err, ok := data.(error); ok {
		data = err.Error()
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Array:
		switch {
		case t == CT_JSON || t == CT_HTML:
			return CT_JSON, data, nil
		case t == CT_XML:
			buff, err := xml.Marshal(data)
			if err != nil {
				panic(err)
				return t, nil, err
			}
			return CT_XML, buff, nil
		case t == CT_YMAL:
			buff, err := yaml.Marshal(data)
			if err != nil {
				return CT_YMAL, nil, err
			}
			return CT_YMAL, buff, nil
		default:
			return t, fmt.Sprintf("%+v", data), nil
		}
	case reflect.String:
		value := []byte(data.(string))
		switch {
		case (t == CT_JSON || t == CT_DEF) && json.Valid(value):
			return CT_JSON, json.RawMessage(value), nil
		case (t == CT_XML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<?xml")):
			return CT_XML, value, nil
		case (t == CT_HTML || t == CT_DEF) && bytes.HasPrefix(value, []byte("<!DOCTYPE html")):
			return CT_HTML, data, nil
		}
		switch {
		case t == CT_JSON:
			return CT_JSON, map[string]interface{}{"data": data}, nil
		default:
			return t, data, nil
		}

	default:
		switch {
		case t == CT_JSON:
			return CT_JSON, map[string]interface{}{"data": data}, nil
		case t == CT_YMAL:
			buff, err := yaml.Marshal(map[string]interface{}{
				"data": data,
			})
			if err != nil {
				return t, nil, err
			}
			return CT_YMAL, buff, nil
		default:
			return t, fmt.Sprint(data), nil
		}
	}
}
