package jsonschema

import (
	"encoding/json"
	"strings"
	"testing"
)

//
//
type A struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"-"`
}

//
//
type Ab struct {
	A string `json:""`
	B string ``
	C string
}

//
//
func TestSimpleProps(t *testing.T) {
	if s := New(A{}); len(s.Props) != 2 {
		t.Error("wrong number of fields in schema")
	} else if len(s.ReqProps) != 0 {
		t.Error("wrong number of required fields in schema")
	}
	if v := New(Ab{}); len(v.Props) != 0 {
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
	if v := New(B{}); len(v.Props) != 1 {
		t.Error("wrong number of fields")
	} else if p, ok := v.Props["b"]; !ok {
		t.Error("missing known schema field")
	} else if mpp, ok := p.(map[string]string); !ok {
		t.Error("failed to cast schema field map")
	} else if mpp["pattern"] != "email" {
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
	s := New(c)

	if d, ok := s.Props["d"]; !ok {
		t.Error("missing known schema field")
	} else if ds, ok := d.(JsonSchema); !ok {
		t.Error("failed to cast nested schema from property")
	} else if mpp, ok := ds.Props["phone"].(map[string]string); !ok {
		t.Error("failed to cast schema field map")
	} else if mpp["pattern"] != "telephone" {
		t.Error("wrong value for pattern key")
	}
}

//
//
type E struct {
	E  string `json:"one"`
	Ea string `json:"two"`
}

//
//
func TestCreateFromStructPointer(t *testing.T) {
	e := E{"eggs", "ham"}

	if s := New(e); len(s.Props) != 2 {
		t.Error("wrong number of fields in schema")
	} else if s := New(&e); len(s.Props) != 2 {
		t.Error("wrong number of fields in schema")
	}
}

//
//
type F struct {
	Inner *F `json:"floop,type=string,format=date-time"`
	Outer *F `json:"gloop,type=string,format=url"`
}

//
//
func TestNestedPointer(t *testing.T) {
	f := &F{Inner: &F{Inner: nil, Outer: nil}, Outer: nil}
	if len(New(f).Props) != 2 {
		t.Error("wrong number of fields in schema")
	}
}

//
//
type G struct {
	G string `json:"g,required"`
	H string `json:"h,required"`
	I string `json:"i"`
}

//
//
func TestMisc(t *testing.T) {
	if s := New(G{}); len(s.ReqProps) != 2 {
		t.Error("wrong number of required fields in schema")
	}
	if s := New(nil); len(s.ReqProps) != 0 || len(s.Props) != 0 {
		t.Error("wrong number of fields on schema")
	}
	i := 0
	if s := New(i); len(s.Props) != 0 {
		t.Error("wrong number of fields on schema")
	}
	if s := New(&i); len(s.Props) != 0 {
		t.Error("wrong number of fields on schema")
	}
	if s := New(""); len(s.Props) != 0 {
		t.Error("wrong number of fields on schema")
	}
}

//
//
func TestEncoding(t *testing.T) {
	s := New(A{})
	if bs, err := json.MarshalIndent(s, "", "  "); err != nil {
		t.Error(err)
	} else if strings.Count(string(bs), "required") != 0 {
		// json schema spec says 'required' key must not be zero length
		t.Error("non-empty required key in output")
	} else {
		t.Logf("encoded a: %v\n", string(bs))
	}

	s = New(G{})
	s.Links = Links{map[string]string{"href": "ok", "rel": "self"}}
	if bs, err := json.MarshalIndent(s, "", "  "); err != nil {
		t.Error(err)
	} else if strings.Count(string(bs), "required") != 1 {
		t.Errorf("wrong number of required keys in output")
	} else {
		t.Logf("encoded b: %v\n", string(bs))
	}
}
