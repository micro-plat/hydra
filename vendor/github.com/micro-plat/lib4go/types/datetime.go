package types

import (
	"encoding/xml"
	"fmt"
	"time"
)

const defaultTimeformat = "2006-01-02 15:04:05"

//DateTime DateTime
type DateTime struct {
	format string
	time.Time
}

//NewDateTime 构建新的DateTime
func NewDateTime(t time.Time, format ...string) *DateTime {
	timefmt := defaultTimeformat
	if len(format) > 0 {
		timefmt = format[0]
	}

	return &DateTime{
		Time:   t,
		format: timefmt,
	}
}

//MarshalJSON MarshalJSON
func (d *DateTime) MarshalJSON() (bytes []byte, err error) {
	val := d.Time
	tmpV := fmt.Sprintf("\"%s\"", val.Format(d.format))
	return []byte(tmpV), nil
}

//UnmarshalJSON UnmarshalJSON
func (d *DateTime) UnmarshalJSON(bytes []byte) error {
	if d.format == "" {
		d.format = defaultTimeformat
	}
	val, err := time.Parse(fmt.Sprintf("\"%s\"", d.format), string(bytes))
	*d = DateTime{Time: val, format: d.format}
	return err
}

//MarshalXML MarshalXML
func (d *DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	val := fmt.Sprintf("%v", d.Time.Format(d.format))
	e.EncodeElement(val, start)
	return nil
}

//UnmarshalXML UnmarshalXML
func (d *DateTime) UnmarshalXML(e *xml.Decoder, start xml.StartElement) error {
	if d.format == "" {
		d.format = defaultTimeformat
	}
	var v string
	e.DecodeElement(&v, &start)
	val, err := time.Parse(d.format, v)
	if err != nil {
		return err
	}
	*d = DateTime{Time: val, format: d.format}
	return nil
}

//Format 默认2006-01-02 15:04:05
func (d *DateTime) Format() string {
	return d.format
}

//String String
func (d *DateTime) String() string {
	return d.Time.Format(d.format)
}
