package vkapi

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// paramsFrom builds request params from struct or map.
// Other types are silently ignored and nil is returned.
//
// Rules:
//
// - param names are detected from "json" struct tags, "-" tag (omit field) and "omitempty" modifiers are respected
//
// - nil keys/values are always omitted (even without "omitempty")
//
// - complex (not string/string slice/number) keys/values are always omitted
//
// - slice values are serialized as comma-separated strings
//
// - bools are serialized as 1 (true) and 0 (false)
//
// REMINDER: keep rules in sync with Client.Exec method doc.
func paramsFrom(d interface{}) url.Values {
	rv := getValue(reflect.ValueOf(d))
	if isEmpty(rv) {
		return nil
	}

	switch rv.Type().Kind() {

	case reflect.Struct:
		return paramsFromStruct(rv)

	case reflect.Map:
		return paramsFromMap(rv)

	default:
		return nil
	}
}

// ----------------------------------------------------------------------------

// paramsFromStruct builds request params from struct.
func paramsFromStruct(rv reflect.Value) url.Values {
	rt := rv.Type()

	q := url.Values{}
	for i, nf := 0, rt.NumField(); i < nf; i++ {
		f := rt.Field(i)
		v := getValue(rv.Field(i))

		if omitValue(v) {
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

// paramsFromMap builds request params from map.
func paramsFromMap(rv reflect.Value) url.Values {
	q := url.Values{}

	for _, k := range rv.MapKeys() {
		k = getValue(k)
		if omitKey(k) {
			continue
		}

		v := getValue(rv.MapIndex(k))
		if omitValue(v) {
			continue
		}

		q.Set(toString(k), toString(v))
	}

	return q
}

// ----------------------------------------------------------------------------

// getValue extracts actual data, dereferencing pointer chain if needed.
// Mostly useful to resolve actual interface value.
func getValue(rv reflect.Value) reflect.Value {
	rv = dereference(rv)
	if rv.Kind() == reflect.Interface {
		rv = getValue(reflect.ValueOf(rv.Interface()))
	}
	return rv
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

// ----------------------------------------------------------------------------

// omitKey tells if key should be omitted (not serialized).
func omitKey(rv reflect.Value) bool {
	switch rv.Kind() {

	case reflect.Slice:
		return true

	default:
		return omitValue(rv)
	}
}

// omitValue tells if value should be omitted (not serialized).
func omitValue(rv reflect.Value) bool {
	vk := rv.Kind()

	// omit unsupported (unserializable) params:
	switch vk {
	case reflect.Invalid,
		reflect.Uintptr,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Struct,
		reflect.UnsafePointer:
		return true
	}

	// omit nil slices:
	if (vk == reflect.Slice || vk == reflect.Array) && rv.IsNil() {
		return true
	}

	// omit invalid params (value of nil pointer etc):
	if !rv.IsValid() {
		return true
	}

	return false
}
