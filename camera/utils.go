package camera

import (
	"encoding/xml"
	"net/http"
	"reflect"
	"strings"
)

func parseSOAPResp(resp *http.Response, dest interface{}) {
	splitedName := strings.Split(reflect.TypeOf(dest).String(), ".")
	typeName := splitedName[len(splitedName)-1]
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
