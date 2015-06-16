package jsonschema

import "testing"

type BasicStruct struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"-"`
}

func TestProps(t *testing.T) {
	b := BasicStruct{}
	s := New("", "", b)
	//t.Logf("%+v %v", s, s.Props)
	if len(*s.Props) != 2 {
		t.Error("wrong number of fields in schema")
	}
}
