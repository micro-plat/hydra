// Copyright 2012-2014 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

// xml.go - basically the core of X2j for map[string]interface{} values.
//          NewMapXml, NewMapXmlReader, mv.Xml, mv.XmlWriter
// see x2j and j2x for wrappers to provide end-to-end transformation of XML and JSON messages.

package anyxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// --------------------------------- Xml, XmlIndent - from mxj -------------------------------

var (
	DefaultRootTag          = "doc"
	UseGoEmptyElementSyntax = false // if 'true' encode empty element as "<tag></tag>" instead of "<tag/>
)

// From: github.com/clbanning/mxj/xml.go with functions relabled: Xml() --> anyxml().
// Encode a Map as XML.  The companion of NewMapXml().
// The following rules apply.
//    - The key label "#text" is treated as the value for a simple element with attributes.
//    - Map keys that begin with a hyphen, '-', are interpreted as attributes.
//      It is an error if the attribute doesn't have a []byte, string, number, or boolean value.
//    - Map value type encoding:
//          > string, bool, float64, int, int32, int64, float32: per "%v" formating
//          > []bool, []uint8: by casting to string
//          > structures, etc.: handed to xml.Marshal() - if there is an error, the element
//            value is "UNKNOWN"
//    - Elements with only attribute values or are null are terminated using "/>".
//    - If len(mv) == 1 and no rootTag is provided, then the map key is used as the root tag, possible.
//      Thus, `{ "key":"value" }` encodes as "<key>value</key>".
//    - To encode empty elements in a syntax consistent with encoding/xml call UseGoXmlEmptyElementSyntax().
func anyxml(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty) // just a stub

	if len(m) == 1 && len(rootTag) == 0 {
		for key, value := range m {
			// if it an array, see if all values are map[string]interface{}
			// we force a new root tag if we'll end up with no key:value in the list
			// so: key:[string_val, bool:true] --> <doc><key>string_val</key><bool>true</bool></key></doc>
			switch value.(type) {
			case []interface{}:
				for _, v := range value.([]interface{}) {
					switch v.(type) {
					case map[string]interface{}: // noop
					default: // anything else
						err = mapToXmlIndent(false, s, DefaultRootTag, m, p)
						goto done
					}
				}
			}
			err = mapToXmlIndent(false, s, key, value, p)
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndent(false, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndent(false, s, DefaultRootTag, m, p)
	}
done:
	return []byte(*s), err
}

// Encode a map[string]interface{} as a pretty XML string.
// See Xml for encoding rules.
func anyxmlIndent(m map[string]interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	if len(m) == 1 && len(rootTag) == 0 {
		// this can extract the key for the single map element
		// use it if it isn't a key for a list
		for key, value := range m {
			if _, ok := value.([]interface{}); ok {
				err = mapToXmlIndent(true, s, DefaultRootTag, m, p)
			} else {
				err = mapToXmlIndent(true, s, key, value, p)
			}
		}
	} else if len(rootTag) == 1 {
		err = mapToXmlIndent(true, s, rootTag[0], m, p)
	} else {
		err = mapToXmlIndent(true, s, DefaultRootTag, m, p)
	}
	return []byte(*s), err
}

type pretty struct {
	indent   string
	cnt      int
	padding  string
	mapDepth int
	start    int
}

func (p *pretty) Indent() {
	p.padding += p.indent
	p.cnt++
}

func (p *pretty) Outdent() {
	if p.cnt > 0 {
		p.padding = p.padding[:len(p.padding)-len(p.indent)]
		p.cnt--
	}
}

var xmlEscapeChars bool

// XMLEscapeChars(true) forces escaping invalid characters in attribute and element values.
// NOTE: this is brute force with NO interrogation of '&' being escaped already; if it is
// then '&amp;' will be re-escaped as '&amp;amp;'.
/*
	The values are:
	"   &quot;
	'   &apos;
	<   &lt;
	>   &gt;
	&   &amp;
*/
func XMLEscapeChars(b bool) {
	xmlEscapeChars = b
}

// order is important - must scan for '&' first
var escapechars = [][2][]byte{
	{[]byte(`&`), []byte(`&amp;`)},
	{[]byte(`<`), []byte(`&lt;`)},
	{[]byte(`>`), []byte(`&gt;`)},
	{[]byte(`"`), []byte(`&quot;`)},
	{[]byte(`'`), []byte(`&apos;`)},
}

func escapeChars(s string) string {
	if len(s) == 0 {
		return s
	}

	b := []byte(s)
	for _, v := range escapechars {
		n := bytes.Count(b, v[0])
		if n == 0 {
			continue
		}
		b = bytes.Replace(b, v[0], v[1], n)
	}
	return string(b)
}

// where the work actually happens
// returns an error if an attribute is not atomic
// patched with new version in github.com/clbanning/mxj - 2015.11.15
func mapToXmlIndent(doIndent bool, s *string, key string, value interface{}, pp *pretty) error {
	var endTag bool
	var isSimple bool
	var elen int
	p := &pretty{pp.indent, pp.cnt, pp.padding, pp.mapDepth, pp.start}

	// per clbanning/mxj issue #48, 18apr18 - try and coerce maps to map[string]interface{}
	if reflect.ValueOf(value).Kind() == reflect.Map {
		switch value.(type) {
		case map[string]interface{}:
		default:
			val := make(map[string]interface{})
			vv := reflect.ValueOf(value)
			keys := vv.MapKeys()
			for _, k := range keys {
				val[fmt.Sprint(k)] = vv.MapIndex(k).Interface()
			}
			value = val
		}
	}

	switch value.(type) {
	case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32:
		if doIndent {
			*s += p.padding
		}
		*s += `<` + key
	}
	switch value.(type) {
	case map[string]interface{}:
		vv := value.(map[string]interface{})
		lenvv := len(vv)
		// scan out attributes - keys have prepended hyphen, '-'
		attrlist := make([][2]string, len(vv))
		var n int
		var ss string
		for k, v := range vv {
			if k[:1] == "-" {
				switch v.(type) {
				case string:
					ss = v.(string)
					if xmlEscapeChars {
						ss = escapeChars(ss)
					}
					attrlist[n][0] = k[1:]
					attrlist[n][1] = ss
				case float64, bool, int, int32, int64, float32:
					attrlist[n][0] = k[1:]
					attrlist[n][1] = fmt.Sprintf("%v", v)
				case []byte:
					ss = string(v.([]byte))
					if xmlEscapeChars {
						ss = escapeChars(ss)
					}
					attrlist[n][0] = k[1:]
					attrlist[n][1] = ss
				default:
					return fmt.Errorf("invalid attribute value for: %s", k)
				}
				n++
			}
		}
		if n > 0 {
			attrlist = attrlist[:n]
			sort.Sort(attrList(attrlist))
			for _, v := range attrlist {
				*s += ` ` + v[0] + `="` + v[1] + `"`
			}
		}

		// only attributes?
		if n == lenvv {
			break
		}
		// simple element? Note: '#text" is an invalid XML tag.
		if v, ok := vv["#text"]; ok {
			if n+1 < lenvv {
				return errors.New("#text key occurs with other non-attribute keys")
			}
			*s += ">" + fmt.Sprintf("%v", v)
			endTag = true
			elen = 1
			isSimple = true
			break
		}
		// close tag with possible attributes
		*s += ">"
		if doIndent {
			*s += "\n"
		}
		// something more complex
		p.mapDepth++
		// extract the map k:v pairs and sort on key
		elemlist := make([][2]interface{}, len(vv))
		n = 0
		for k, v := range vv {
			if k[:1] == "-" {
				continue
			}
			elemlist[n][0] = k
			elemlist[n][1] = v
			n++
		}
		elemlist = elemlist[:n]
		sort.Sort(elemList(elemlist))
		var i int
		for _, v := range elemlist {
			switch v[1].(type) {
			case []interface{}:
			default:
				if i == 0 && doIndent {
					p.Indent()
				}
			}
			i++
			mapToXmlIndent(doIndent, s, v[0].(string), v[1], p)
			switch v[1].(type) {
			case []interface{}: // handled in []interface{} case
			default:
				if doIndent {
					p.Outdent()
				}
			}
			i--
		}
		p.mapDepth--
		endTag = true
		elen = 1 // we do have some content ...
	case []interface{}:
		for _, v := range value.([]interface{}) {
			if doIndent {
				p.Indent()
			}
			mapToXmlIndent(doIndent, s, key, v, p)
			if doIndent {
				p.Outdent()
			}
		}
		return nil
	case nil:
		// terminate the tag
		if doIndent {
			*s += p.padding
		}
		*s += "<" + key
		endTag, isSimple = true, true
		break
	default: // handle anything - even goofy stuff
		elen = 0
		switch value.(type) {
		case string:
			v := value.(string)
			if xmlEscapeChars {
				v = escapeChars(v)
			}
			elen = len(v)
			if elen > 0 {
				*s += ">" + v
			}
		case float64, bool, int, int32, int64, float32:
			v := fmt.Sprintf("%v", value)
			elen = len(v) // always > 0
			*s += ">" + v
		case []byte: // NOTE: byte is just an alias for uint8
			// similar to how xml.Marshal handles []byte structure members
			v := string(value.([]byte))
			if xmlEscapeChars {
				v = escapeChars(v)
			}
			elen = len(v)
			if elen > 0 {
				*s += ">" + v
			}
		default:
			var v []byte
			var err error
			if doIndent {
				v, err = xml.MarshalIndent(value, p.padding, p.indent)
			} else {
				v, err = xml.Marshal(value)
			}
			if err != nil {
				*s += ">UNKNOWN"
			} else {
				elen = len(v)
				if elen > 0 {
					*s += string(v)
				}
			}
		}
		isSimple = true
		endTag = true
	}
	if endTag {
		if doIndent {
			if !isSimple {
				*s += p.padding
			}
		}
		switch value.(type) {
		case map[string]interface{}, []byte, string, float64, bool, int, int32, int64, float32, nil:
			if elen > 0 || UseGoEmptyElementSyntax {
				if elen == 0 {
					*s += ">"
				}
				*s += `</` + key + ">"
			} else {
				*s += `/>`
			}
		}
	} else if UseGoEmptyElementSyntax {
		*s += "></" + key + ">"
	} else {
		*s += "/>"
	}
	if doIndent {
		if p.cnt > p.start {
			*s += "\n"
		}
		p.Outdent()
	}

	return nil
}

// ============================ sort interface implementation =================

type attrList [][2]string

func (a attrList) Len() int {
	return len(a)
}

func (a attrList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a attrList) Less(i, j int) bool {
	if a[i][0] > a[j][0] {
		return false
	}
	return true
}

type elemList [][2]interface{}

func (e elemList) Len() int {
	return len(e)
}

func (e elemList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e elemList) Less(i, j int) bool {
	if e[i][0].(string) > e[j][0].(string) {
		return false
	}
	return true
}
