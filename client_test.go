package vkapi_test

import (
	"net/http"
	"net/url"

	"github.com/mxmCherry/vkapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var httpClient *mockHTTPClient
	var subject *vkapi.Client

	BeforeEach(func() {
		httpClient = &mockHTTPClient{}
		subject = vkapi.From(httpClient, vkapi.Options{
			AccessToken: "DUMMY_TOKEN",
			Version:     "42.42",
		})
	})

	Describe("Exec", func() {

		BeforeEach(func() {
			httpClient.code = http.StatusOK
			httpClient.body = `{}`
		})

		It("should build request URL", func() {
			err := subject.Exec("dummy.method", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(httpClient.url).To(Equal("https://api.vk.com/method/dummy.method"))
		})

		It("should use provided params", func() {
			params := url.Values{
				"param_name": []string{"PARAM_VALUE"},
			}

			err := subject.Exec("dummy.users.search", params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.form.Get("param_name")).To(Equal("PARAM_VALUE"))
		})

		It("should parse response", func() {
			httpClient.body = `{
				"response": {
					"count": 1,
					"items": [
						{
							"id": 42,
							"first_name": "FirstName",
							"last_name": "LastName"
						}
					]
				}
			}`

			response := new(struct {
				Count uint64 `json:"count"`
				Items []struct {
					ID        uint64 `json:"id"`
					FirstName string `json:"first_name"`
					LastName  string `json:"last_name"`
				} `json:"items"`
			})

			err := subject.Exec("dummy.method", nil, response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Count).To(Equal(uint64(1)))
			Expect(response.Items).To(HaveLen(1))
			Expect(response.Items[0].ID).To(Equal(uint64(42)))
			Expect(response.Items[0].FirstName).To(Equal("FirstName"))
			Expect(response.Items[0].LastName).To(Equal("LastName"))
		})

		It("should return vk errors", func() {
			httpClient.code = http.StatusOK
			httpClient.body = `{
				"error": {
					"error_code": 42,
					"error_msg": "Test error"
				}
			}`

			err := subject.Exec("", nil, nil)
			Expect(err).To(Equal(vkapi.Error{
				ErrorCode: 42,
				ErrorMsg:  "Test error",
			}))
		})

		Context("access token", func() {

			It("should use client access token", func() {
				err := subject.Exec("dummy.method", nil, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("access_token"), "DUMMY_TOKEN")
			})

			It("should allow to override access token for a single request", func() {
				params := url.Values{
					"access_token": []string{"OVERRIDDEN_TOKEN"},
				}

				err := subject.Exec("dummy.method", params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("access_token")).To(Equal("OVERRIDDEN_TOKEN"))
			})

			It("should allow to clear access token for a single request", func() {
				params := url.Values{
					"access_token": []string{},
				}

				err := subject.Exec("dummy.method", params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("access_token")).To(Equal(""))
			})

		}) // access token context

		Context("version", func() {

			It("should use client version", func() {
				err := subject.Exec("dummy.method", nil, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("v"), "42.42")
			})

			It("should allow to override version for a single request", func() {
				params := url.Values{
					"v": []string{"99.99"},
				}

				err := subject.Exec("dummy.method", params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("v")).To(Equal("99.99"))
			})

			It("should allow to clear version for a single request", func() {
				params := url.Values{
					"v": []string{},
				}

				err := subject.Exec("dummy.method", params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(httpClient.form.Get("v")).To(Equal(""))
			})

		}) // version context

	}) // Exec description

})
