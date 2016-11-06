package vkapi

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ToParams builds request params from struct.
// It returns nil params unless arg is not a struct.
//
// Param names are detected from "json" struct tags,
// "-" tag (omit field) and "omitempty" modifier are respected.
//
// Nil values are always omitted (even without "omitempty").
// Chan, func, interface, map and complex values are also omitted.
//
// Slices are serialized as comma-separated strings.
//
// Bools are serialized as 1 (true) and 0 (false).
func ToParams(d interface{}) url.Values {
	rv := dereference(reflect.ValueOf(d))
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
