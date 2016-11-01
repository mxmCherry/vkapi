package vkapi_test

import (
	"net/url"

	"github.com/mxmCherry/vkapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ToParams", func() {
	It("should build params from struct", func() {
		d := struct {
			String   string   `json:"string_param"`
			Number   uint64   `json:"number_param"`
			Repeated []uint64 `json:"repeated_param"`
			Empty    string   `json:"empty,omitempty"`
		}{
			String:   "string_value",
			Number:   42,
			Repeated: []uint64{1, 2, 3},
			Empty:    "",
		}
		Expect(vkapi.ToParams(d)).To(Equal(url.Values{
			"string_param":   []string{"string_value"},
			"number_param":   []string{"42"},
			"repeated_param": []string{"1,2,3"},
		}))
	})
})
