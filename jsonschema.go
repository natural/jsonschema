package jsonschema

import (
	"reflect"
	"strings"
)

//
//
type Props map[string]interface{}

//
//
type ReqProps []string

//
//
type JsonSchema struct {
	Schema   string   `json:"$schema,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Desc     string   `json:"description,omitempty"`
	AddProps bool     `json:"additionalProperties,omitempty"`
	ReqProps ReqProps `json:"required,omitempty"`
	Props    Props    `json:"properties,omitempty"`
}

//
//
func New(name, desc string, v interface{}) JsonSchema {
	rp, p := props(v)
	return JsonSchema{
		Schema:   "http://json-schema.org/schema#",
		Name:     name,
		Type:     "object",
		Desc:     desc,
		ReqProps: rp,
		Props:    p,
	}
}

//
//
func props(v interface{}) (ReqProps, Props) {
	pm := map[string]interface{}{}
	pr := []string{}

	for _, field := range fields("json", v) {
		js, ft := field.Tag.Get("json"), field.Type
		if js == "" || js == "-" {
			continue
		}
		nm, ps, rs := parsetag(js)
		if ft.Kind() == reflect.Struct {
			vv := reflect.Indirect(reflect.New(ft)).Interface()
			nv := New(nm, "", vv)
			nv.Schema = ""
			pm[nm] = nv
		} else {
			pm[nm] = ps
		}
		for _, v := range rs {
			pr = append(pr, v)
		}
		if v, n := types[ft.Kind()], ps["type"]; v != "" && n == "" {
			ps["type"] = v
		}
	}
	return pr, pm
}

//
//
func fields(name string, src interface{}) []reflect.StructField {
	fs := []reflect.StructField{}
	st := reflect.TypeOf(src)
	for i := 0; i < st.NumField(); i++ {
		if n := st.Field(i).Tag.Get(name); n != "" {
			fs = append(fs, st.Field(i))
		}
	}
	return fs
}

//
//
func parsetag(v string) (string, map[string]string, []string) {
	vs := map[string]string{}
	rs := []string{}
	nm := ""

	for i, s := range strings.Split(v, ",") {
		if i == 0 {
			nm = s
			continue
		}
		sp := strings.Split(s, "=")
		if c := len(sp); c == 1 {
			rs = append(rs, nm)
		} else if c == 2 {
			vs[sp[0]] = sp[1]
		}
	}
	return nm, vs, rs
}

//
//
var (
	types = map[reflect.Kind]string{
		reflect.Bool:          "boolean",
		reflect.Int:           "integer",
		reflect.Int8:          "integer",
		reflect.Int16:         "integer",
		reflect.Int32:         "integer",
		reflect.Int64:         "integer",
		reflect.Uint:          "integer",
		reflect.Uint8:         "integer",
		reflect.Uint16:        "integer",
		reflect.Uint32:        "integer",
		reflect.Uint64:        "integer",
		reflect.Uintptr:       "null",
		reflect.Float32:       "number",
		reflect.Float64:       "number",
		reflect.Complex64:     "number",
		reflect.Complex128:    "number",
		reflect.Array:         "array",
		reflect.Chan:          "object",
		reflect.Func:          "object",
		reflect.Interface:     "object",
		reflect.Map:           "object",
		reflect.Ptr:           "object",
		reflect.Slice:         "array",
		reflect.String:        "string",
		reflect.Struct:        "object",
		reflect.UnsafePointer: "null",
	}
)
