package vkapi

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("paramsFrom", func() {

	It("should build params from struct", func() {
		s := "string"
		d := struct {
			Bool     bool     `json:"bool_param"`
			String   string   `json:"string_param"`
			Number   uint64   `json:"number_param"`
			Repeated []uint64 `json:"repeated_param"`
			Empty    string   `json:"empty_param,omitempty"`
			Omitted  string   `json:"-"`
			Pointer  *string  `json:"pointer_param"`
			Nil      *string  `json:"nil_param"`
		}{
			Bool:     true,
			String:   "string_value",
			Number:   42,
			Repeated: []uint64{1, 2, 3},
			Empty:    "",
			Omitted:  "omitted",
			Pointer:  &s,
			Nil:      nil,
		}
		Expect(paramsFrom(d)).To(Equal(url.Values{
			"bool_param":     []string{"1"},
			"string_param":   []string{"string_value"},
			"number_param":   []string{"42"},
			"repeated_param": []string{"1,2,3"},
			"pointer_param":  []string{"string"},
		}))
	})

	It("should build params from map", func() {
		s := "string"
		d := map[string]interface{}{
			"bool_param":     true,
			"string_param":   "string_value",
			"number_param":   42,
			"repeated_param": []uint64{1, 2, 3},
			"pointer_param":  &s,
			"nil_param":      nil,
		}
		Expect(paramsFrom(d)).To(Equal(url.Values{
			"bool_param":     []string{"1"},
			"string_param":   []string{"string_value"},
			"number_param":   []string{"42"},
			"repeated_param": []string{"1,2,3"},
			"pointer_param":  []string{"string"},
		}))
	})

	It("should ignore unsupported types", func() {
		Expect(paramsFrom(nil)).To(BeNil())
		Expect(paramsFrom(42)).To(BeNil())
		Expect(paramsFrom("string")).To(BeNil())
	})

})
