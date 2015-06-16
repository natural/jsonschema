package jsonschema

import (
	"reflect"
	"strings"
)

//
//
func New(name, desc string, v interface{}) map[string]interface{} {
	rp, ap := Properties(v)
	return map[string]interface{}{
		"$schema":              "http://json-schema.org/schema#",
		"name":                 name,
		"type":                 "object",
		"description":          desc,
		"additionalProperties": false,

		"required":   rp,
		"properties": ap,
	}
}

//
//
func Properties(v interface{}) ([]string, map[string]interface{}) {
	pm := map[string]interface{}{}
	rp := []string{}

	for _, field := range Fields("json", v) {
		js, ft := field.Tag.Get("json"), field.Type
		nm, ps := ParseTagValue(js)
		pm[nm] = ps
		if v, n := FieldKindTypeMap[ft.Kind()], ps["type"]; v != "" && n == "" {
			ps["type"] = v
		}
	}
	return rp, pm
}

func ParseTagValue(v string) (string, map[string]string) {
	vs := map[string]string{}
	nm := ""

	for i, s := range strings.Split(v, ",") {
		if i == 0 {
			nm = s
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
	FieldNameTypeMap = map[string]string{
		"Time": "string",
	}

	FieldKindTypeMap = map[reflect.Kind]string{
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

func Fields(name string, src interface{}) []reflect.StructField {
	fs := []reflect.StructField{}
	st := reflect.TypeOf(src)
	for i := 0; i < st.NumField(); i++ {
		if n := st.Field(i).Tag.Get(name); n != "" {
			fs = append(fs, st.Field(i))
		}
	}
	return fs
}
