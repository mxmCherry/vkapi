package vkapi

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var subject Client

	BeforeEach(func() {
		subject = New(Options{
			AccessToken: "dummy_token",
		})
		subject.(*client).http = &mockHttpClient{
			code: http.StatusNotFound,
			body: `{
				"error": {
					"error_code": 0,
					"error_msg": "Not Found"
				}
			}`,
		}
	})

	It("should exec", func() {
		subject.(*client).http.(*mockHttpClient).code = http.StatusOK
		subject.(*client).http.(*mockHttpClient).body = `{
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

		query := url.Values{
			"q": []string{"FirstName LastName"},
		}

		err := subject.Exec("dummy.users.search", query, response)
		Expect(err).NotTo(HaveOccurred())

		Expect(subject.(*client).http.(*mockHttpClient).url).To(Equal(
			"https://api.vk.com/method/dummy.users.search?access_token=dummy_token&q=FirstName+LastName&v=5.59",
		))

		Expect(response.Count).To(Equal(uint64(1)))
		Expect(response.Items).To(HaveLen(1))
		Expect(response.Items[0].ID).To(Equal(uint64(42)))
		Expect(response.Items[0].FirstName).To(Equal("FirstName"))
		Expect(response.Items[0].LastName).To(Equal("LastName"))
	})

	It("should return vk errors", func() {
		subject.(*client).http.(*mockHttpClient).code = http.StatusOK
		subject.(*client).http.(*mockHttpClient).body = `{
			"error": {
				"error_code": 42,
				"error_msg": "ErrorMsg"
			}
		}`

		err := subject.Exec("", nil, nil)
		Expect(err).To(MatchError("vkapi: ErrorMsg (code 42)"))
	})

})
