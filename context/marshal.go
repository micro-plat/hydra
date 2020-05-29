package context

import (
	"encoding/xml"
	"io"
)

//UnmarshalXML 将xml反序列化为map
func UnmarshalXML(s string) (map[string]string, error) {
	output := make(map[string]string)
	err := xml.Unmarshal([]byte(s), (*xmlMap)(&output))
	return output, err
}

type xmlMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m xmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *xmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = xmlMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
