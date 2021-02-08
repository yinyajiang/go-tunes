package mtunes

import (
	"fmt"
	"reflect"

	"howett.net/plist"
)

//UnmashalPlist ...
func UnmashalPlist(data []byte, v interface{}) (format string, err error) {
	if !reflect.ValueOf(v).Elem().CanAddr() {
		err = fmt.Errorf("Must input a pointer")
		return
	}
	if len(data) == 0 {
		err = fmt.Errorf("Data is empty")
		return
	}

	nformat, e := plist.Unmarshal(data, v)
	if e != nil {
		err = e
		return
	}
	format = plist.FormatNames[nformat]
	return
}

//MashalPlist ...
func MashalPlist(v interface{}) ([]byte, error) {
	return plist.Marshal(v, plist.BinaryFormat)
}

//MashalPlistString ...
func MashalPlistString(v interface{}) (s string, err error) {
	r, err := plist.Marshal(v, plist.XMLFormat)
	s = string(r)
	return
}

//PlistToString ...
func PlistToString(data []byte) string {
	var dict map[string]interface{}
	format, err := UnmashalPlist(data, &dict)
	if format == "XML" {
		return string(data)
	}
	data, err = plist.Marshal(dict, plist.XMLFormat)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}
