package vkapi

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// valuesFromStruct builds request params from struct.
// Other types are silently ignored and nil is returned.
//
// Structs are interpreted using the following rules:
//
// - param names are detected from "json" struct tags, "-" tag (omit field) and "omitempty" modifiers are respected
//
// - nil values are always omitted (even without "omitempty")
//
// - chan, func, interface, map, struct and complex values are always omitted
//
// - slices are serialized as comma-separated strings
//
// - bools are serialized as 1 (true) and 0 (false)
//
// REMINDER: keep this in sync with Client.Exec method doc.
func valuesFromStruct(d interface{}) url.Values {
	rv := dereference(reflect.ValueOf(d))
	if isEmpty(rv) {
		return nil
	}

	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return nil
	}

	q := url.Values{}
	for i, nf := 0, rt.NumField(); i < nf; i++ {
		f := rt.Field(i)
		v := dereference(rv.Field(i))
		vk := v.Kind()

		// omit unsupported (unserializable) params:
		switch vk {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Uintptr, reflect.UnsafePointer,
			reflect.Complex64, reflect.Complex128:
			continue
		}

		// omit nil slices:
		if vk == reflect.Slice && v.IsNil() {
			continue
		}

		// omit invalid params (value of nil pointer etc):
		if !v.IsValid() {
			continue
		}

		name := ""
		omitempty := false

		tag := strings.Split(f.Tag.Get("json"), ",")
		if n := len(tag); n >= 1 {
			name = tag[0]
			if name == "-" {
				continue
			}
			if n >= 2 {
				for _, s := range tag {
					if s == "omitempty" {
						omitempty = true
						break
					}
				}
			}
		} else {
			name = strings.ToLower(f.Name)
		}

		if omitempty && isEmpty(v) {
			continue
		}

		q.Set(name, toString(v))
	}

	return q
}

// dereference dereferences pointer chain.
func dereference(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	return rv
}

// toString stringifies value.
func toString(rv reflect.Value) string {
	rv = dereference(rv)
	if isEmpty(rv) {
		return ""
	}

	vi := rv.Interface()
	switch x := vi.(type) {

	case bool:
		if x {
			return "1"
		}
		return "0"

	case []byte:
		return string(x)

	default:
		if rv.Kind() == reflect.Slice {
			n := rv.Len()
			ss := make([]string, n)
			for j := 0; j < n; j++ {
				ss[j] = toString(dereference(rv.Index(j)))
			}
			return strings.Join(ss, ",")
		}
		return fmt.Sprintf("%v", vi)

	}
}

// isEmpty tells if value is empty.
func isEmpty(rv reflect.Value) bool {
	if !rv.IsValid() {
		return true
	}

	zero := reflect.Zero(rv.Type())
	if !zero.IsValid() {
		return false
	}

	return reflect.DeepEqual(rv.Interface(), zero.Interface())
}
