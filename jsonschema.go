package jsonschema

import (
	"reflect"
	"strings"
)

//
//
type JsonSchema struct {
	Schema   string    `json:"$schema"`
	Name     string    `json:"name,omitempty"`
	Type     string    `json:"type,omitempty"`
	Desc     string    `json:"description,omitempty"`
	AddProp  bool      `json:"additionalProperties,omitempty"`
	ReqProps *[]string `json:"required,omitempty"`

	Props *map[string]interface{} `json:"properties,omitempty"`
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
func props(v interface{}) (*[]string, *map[string]interface{}) {
	pm := map[string]interface{}{}
	pr := &[]string{}

	for _, field := range fields("json", v) {
		js, ft := field.Tag.Get("json"), field.Type
		if js == "" || js == "-" {
			continue
		}
		nm, ps := parsetag(js)
		pm[nm] = ps
		if v, n := types[ft.Kind()], ps["type"]; v != "" && n == "" {
			ps["type"] = v
		}
	}
	if len(*pr) == 0 {
		pr = nil
	}
	ppm := &pm
	if len(pm) == 0 {
		ppm = nil
	}
	return pr, ppm
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
func parsetag(v string) (string, map[string]string) {
	vs := map[string]string{}
	nm := ""

	for i, s := range strings.Split(v, ",") {
		if i == 0 {
			nm = s
			continue
		}
		if s == "omitempty" {
			continue
		}
		sp := strings.Split(s, "=")
		if len(sp) == 2 {
			vs[sp[0]] = sp[1]
		}
	}
	return nm, vs
}

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
