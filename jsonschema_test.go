package jsonschema

import "testing"

//
//
type A struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"-"`
}

//
//
func TestSimpleProps(t *testing.T) {
	v := New("", "", A{})
	if len(v.Props) != 2 {
		t.Error("wrong number of fields in schema")
	}
}

//
type Ab struct {
	A string `json:""`
	B string ``
	C string
}

//
//
func TestEmptyProps(t *testing.T) {
	v := New("", "", Ab{})
	if len(v.Props) != 0 {
		t.Error("wrong number of fields in schema")
	}
}

//
//
type B struct {
	B int `json:"b,pattern=email"`
}

//
//
func TestKeywordProps(t *testing.T) {
	v := New("", "", B{})
	if len(v.Props) != 1 {
		t.Error("wrong number of fields")
	}
	p, ok := v.Props["b"]
	if !ok {
		t.Error("missing known schema field")
	}
	mpp, ok := p.(map[string]string)
	if !ok {
		t.Error("failed to cast schema field map")
	}
	if mpp["pattern"] != "email" {
		t.Error("wrong value for pattern key")
	}
}

//
//
type C struct {
	C string `json:"c,pattern=url"`
	D D      `json:"d"`
}

//
//
type D struct {
	phone string `json:"phone,pattern=telephone"`
}

//
//
func TestNestedProps(t *testing.T) {
	c := C{C: "anything", D: D{phone: "again"}}
	s := New("type-name-c", "type-desc-c", c)
	d, ok := s.Props["d"]
	if !ok {
		t.Error("missing known schema field")
	}
	ds, ok := d.(JsonSchema)
	if !ok {
		t.Error("failed to cast nested schema from property")
	}
	mpp, ok := ds.Props["phone"].(map[string]string)
	if !ok {
		t.Error("failed to cast schema field map")
	}
	if mpp["pattern"] != "telephone" {
		t.Error("wrong value for pattern key")
	}
}
