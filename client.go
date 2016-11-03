package vkapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// Client abstracts vk.com API client: https://vk.com/dev/api_requests
type Client interface {
	// Exec calls vk.com API method: https://vk.com/dev/methods
	//
	// Response arg must be a pointer to unmarshal "response" field:
	//   {
	//     "response": {this data will be unmarshalled into response arg}
	//   }
	Exec(method string, params url.Values, response interface{}) error
}

// Options hold configuration data for Client.
type Options struct {
	// AccessToken holds vk.com API access token (optional): https://vk.com/dev/access_token
	AccessToken string
}

// HTTPClient abstracts HTTP client.
type HTTPClient interface {
	PostForm(string, url.Values) (*http.Response, error)
}

// New creates new Client with default HTTP client.
func New(options Options) Client {
	return From(new(http.Client), options)
}

// From creates new Client from custom (preconfigured) HTTP client.
// It may be used, for example, if proxy support is needed.
func From(httpClient HTTPClient, options Options) Client {
	return &client{
		options: options,
		http:    httpClient,
	}
}

// ----------------------------------------------------------------------------

const (
	scheme  = "https"
	host    = "api.vk.com"
	prefix  = "/method/"
	version = "5.59"
)

type client struct {
	options Options
	http    HTTPClient
}

func (c *client) Exec(method string, params url.Values, response interface{}) error {
	if params == nil {
		params = url.Values{}
	}
	if c.options.AccessToken != "" {
		params.Set("access_token", c.options.AccessToken)
	}
	params.Set("v", version)

	resp, err := c.http.PostForm((&url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path.Join(prefix, method),
	}).String(), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wrapper := &struct {
		Error    *Error      `json:"error,omitempty"`
		Response interface{} `json:"response"`
	}{
		Response: response,
	}
	if err := json.NewDecoder(resp.Body).Decode(wrapper); err != nil {
		return err
	}

	if wrapper.Error != nil {
		return *wrapper.Error
	}
	return nil
}
