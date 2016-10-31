package vkapi

import (
	"net/http"
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

type mockHttpClient struct {
	code int
	body string
	err  error
	url  string
}

func (c *mockHttpClient) Get(url string) (*http.Response, error) {
	c.url = url
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
