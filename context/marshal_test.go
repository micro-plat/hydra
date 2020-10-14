package context

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestUnmarshalXML(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{name: "转换正确的xml", args: args{s: `<xml><key>value</key><key1>value1</key1></xml>`}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalXML(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalXML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_xmlMap_UnmarshalXML(t *testing.T) {
	type args struct {
		d     *xml.Decoder
		start xml.StartElement
	}
	tests := []struct {
		name    string
		m       *xmlMap
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalXML(tt.args.d, tt.args.start); (err != nil) != tt.wantErr {
				t.Errorf("xmlMap.UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
