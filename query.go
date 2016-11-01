package vkapi

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Query builds query from struct.
// It returns nil query unless arg is not a struct.
// Param names are detected from "json" struct tags.
// Slices are serialized as comma-separated strings.
func Query(d interface{}) url.Values {
	rv := reflect.ValueOf(d)
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return nil
	}

	q := url.Values{}
	for i, nf := 0, rt.NumField(); i < nf; i++ {
		f := rt.Field(i)
		v := rv.Field(i)
		vi := v.Interface()

		name := strings.ToLower(f.Name)
		omitempty := false

		tag := strings.Split(f.Tag.Get("json"), ",")
		if n := len(tag); n >= 1 {
			name = tag[0]
			if n >= 2 {
				for _, s := range tag {
					if s == "omitempty" {
						omitempty = true
						break
					}
				}
			}
		}

		if omitempty && reflect.DeepEqual(vi, reflect.Zero(f.Type).Interface()) {
			continue
		}

		if b, ok := vi.([]byte); ok {
			vi = string(b)
		} else if f.Type.Kind() == reflect.Slice {
			n := v.Len()
			ss := make([]string, n)
			for j := 0; j < n; j++ {
				ss[j] = fmt.Sprintf("%v", v.Index(j).Interface())
			}
			vi = strings.Join(ss, ",")
		}

		q.Set(name, fmt.Sprintf("%v", vi))
	}

	return q
}
