// Package jsonschema provides a simple approach for deriving JSON Schemas
// from user-defined structs.
//
package jsonschema

import (
	"reflect"
	"strings"
)

// Props maps property keys to arbitrary values.
//
type Props map[string]interface{}

// ReqProps indicates required property keys.
//
type ReqProps []string

// Links is a slice of anything (until a better approach is found).
//
type Links []interface{}

// JsonSchema encapsulates the fields of a serializable JSON schema.
//
type JsonSchema struct {
	Schema   string   `json:"$schema,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Desc     string   `json:"description,omitempty"`
	AddProps bool     `json:"additionalProperties,omitempty"`
	ReqProps ReqProps `json:"required,omitempty"`
	Props    Props    `json:"properties,omitempty"`
	Links    Links    `json:"links,omitempty"`
}

// New creates and returns a JsonSchema from the given value (struct).
// The first two optional arguments are interpreted as Name and Description.
//
func New(v interface{}, opts ...string) JsonSchema {
	nm, ds := "", ""
	c := len(opts)
	if c > 0 {
		nm = opts[0]
	}
	if c > 1 {
		ds = opts[1]
	}
	rp, p := props(v)
	return JsonSchema{
		Schema:   "http://json-schema.org/schema#",
		Name:     nm,
		Type:     "object",
		Desc:     ds,
		ReqProps: rp,
		Props:    p,
	}
}

// This derives the required properties slice and the properties map from
// the given value.
//
func props(v interface{}) (ReqProps, Props) {
	pr, pm := []string{}, map[string]interface{}{}

	for _, f := range fields("json", v) {
		js, ft := f.Tag.Get("json"), f.Type
		if js == "" || js == "-" {
			continue
		}
		nm, ps, rs := parsetag(js)
		for _, v := range rs {
			pr = append(pr, v)
		}
		if ft.Kind() == reflect.Struct && ps["type"] == "" {
			vv := reflect.Indirect(reflect.New(ft)).Interface()
			nv := New(vv)
			nv.Name = nm
			nv.Schema = ""
			pm[nm] = nv
		} else {
			pm[nm] = ps
		}
		if v, n := types[ft.Kind()], ps["type"]; v != "" && n == "" {
			ps["type"] = v
		}
	}
	// who says go doesn't do automatic type conversion?
	return pr, pm
}

// This returns the struct fields that have the given struct tag.
//
func fields(name string, src interface{}) []reflect.StructField {
	fs := []reflect.StructField{}
	st := reflect.TypeOf(src)
	if st == nil {
		return fs
	} else if k := st.Kind(); k == reflect.Ptr {
		st = st.Elem()
	}
	// retest
	if st.Kind() != reflect.Struct {
		return fs
	}
	//fmt.Printf("KIND: %v %T\n", st, st)
	for i := 0; i < st.NumField(); i++ {
		if n := st.Field(i).Tag.Get(name); n != "" {
			fs = append(fs, st.Field(i))
		}
	}
	return fs
}

// "Parse" is used loosely here.g
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
		reflect.Complex64:     "object",
		reflect.Complex128:    "object",
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
