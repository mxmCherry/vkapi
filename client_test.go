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
	var subject vkapi.Client

	BeforeEach(func() {
		httpClient = &mockHTTPClient{
			code: http.StatusNotFound,
			body: `{
				"error": {
					"error_code": 0,
					"error_msg": "Not Found"
				}
			}`,
		}
		subject = vkapi.From(httpClient, vkapi.Options{
			AccessToken: "dummy_token",
		})
	})

	It("should exec", func() {
		httpClient.code = http.StatusOK
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

		query := url.Values{
			"q": []string{"FirstName LastName"},
		}

		response := new(struct {
			Count uint64 `json:"count"`
			Items []struct {
				ID        uint64 `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"items"`
		})

		err := subject.Exec("dummy.users.search", query, response)
		Expect(err).NotTo(HaveOccurred())

		Expect(httpClient.url).To(Equal(
			"https://api.vk.com/method/dummy.users.search?access_token=dummy_token&q=FirstName+LastName&v=5.59",
		))

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
				"error_msg": "ErrorMsg"
			}
		}`

		err := subject.Exec("", nil, nil)
		Expect(err).To(MatchError("vkapi: ErrorMsg (code 42)"))
	})

})