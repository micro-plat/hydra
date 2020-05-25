package dispatcher

import (
	"encoding/xml"
	"path"
	"reflect"
	"runtime"
)

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// H is a shortcup for map[string]interface{}
type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
