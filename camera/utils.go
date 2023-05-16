package camera

import (
	"encoding/xml"
	"net/http"
	"reflect"
)

func parseSOAPResp(resp *http.Response, dest interface{}) {
	typeName := reflect.TypeOf(dest).Name()
	d := xml.NewDecoder(resp.Body)
	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
		switch start := t.(type) {
		case xml.StartElement:
			if start.Name.Local == typeName {
				d.DecodeElement(dest, &start)
				return
			}
		}
	}
}
