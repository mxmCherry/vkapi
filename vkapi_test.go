package vkapi_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVkApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "vkapi")
}

// ----------------------------------------------------------------------------

type mockHTTPClient struct {
	code int
	body string
	err  error
	url  string
	form url.Values
}

func (c *mockHTTPClient) PostForm(url string, form url.Values) (*http.Response, error) {
	c.url = url
	c.form = form
	return &http.Response{
		StatusCode: c.code,
		Body: mockReadCloser{
			Reader: strings.NewReader(c.body),
		},
	}, c.err
}

type mockReadCloser struct{ *strings.Reader }

func (rc mockReadCloser) Close() error {
	return nil
}
