// anyxml - marshal an XML document from almost any Go variable.
// Marshal XML from map[string]interface{}, arrays, slices, alpha/numeric, etc.
//
// Wraps xml.Marshal with functionality in github.com/clbanning/mxj to create
// a more genericized XML marshaling capability. Note: unmarshaling the resultant
// XML may not return the original value, since tag labels may have been injected
// to create the XML representation of the value.
//
// See mxj package documentation for more information.  See anyxml_test.go for
// examples or just try Xml() or XmlIndent().
/*
 Encode an arbitrary JSON object.
	package main

	import (
		"encoding/json"
		"fmt"
		"github.com/clbanning/anyxml"
	)

	func main() {
		jsondata := []byte(`[
			{ "somekey":"somevalue" },
			"string",
			3.14159265,
			true
		]`)
		var i interface{}
		err := json.Unmarshal(jsondata, &i)
		if err != nil {
			// do something
		}
		x, err := anyxml.XmlIndent(i, "", "  ", "mydoc")
		if err != nil {
			// do something else
		}
		fmt.Println(string(x))
	}

	output:
		<mydoc>
		  <somekey>somevalue</somekey>
		  <element>string</element>
		  <element>3.14159265</element>
		  <element>true</element>
		</mydoc>
*/
package anyxml

import (
	"encoding/xml"
	"reflect"
)

// Encode arbitrary value as XML.  Note: there are no guarantees.
func Xml(v interface{}, rootTag ...string) ([]byte, error) {
	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}
	if v == nil {
		if UseGoEmptyElementSyntax {
			return []byte("<" + rt + "></" + rt + ">"), nil
		}
		return []byte("<" + rt + "/>"), nil
	}
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.Marshal(v)
	}

	var err error
	s := new(string)
	p := new(pretty)

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">"
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndent(false, s, tag, val, p)
					}
				} else {
					err = mapToXmlIndent(false, s, "element", vv, p)
				}
			default:
				err = mapToXmlIndent(false, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		b, err = anyxml(v.(map[string]interface{}), rootTag...)
	default:
		err = mapToXmlIndent(false, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

// Encode an arbitrary value as a pretty XML string. Note: there are no guarantees.
func XmlIndent(v interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}
	if v == nil {
		if UseGoEmptyElementSyntax {
			return []byte(prefix + "<" + rt + ">\n" + prefix + "</" + rt + ">"), nil
		}
		return []byte(prefix + "<" + rt + "/>"), nil
	}
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
	}

	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">\n"
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndent(true, s, tag, val, p)
					}
				} else {
					p.start = 1 // we're 1 tag in to the doc
					err = mapToXmlIndent(true, s, "element", vv, p)
					*s += "\n"
				}
			default:
				p.start = 0
				err = mapToXmlIndent(true, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		b, err = anyxmlIndent(v.(map[string]interface{}), prefix, indent, rootTag...)
	default:
		err = mapToXmlIndent(true, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}
