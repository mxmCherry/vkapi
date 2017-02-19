package vkapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// Client represents vk.com API client:
// https://vk.com/dev/api_requests
type Client interface {
	// Exec calls vk.com API method:
	// https://vk.com/dev/methods
	//
	// Response arg must be a pointer to unmarshal "response" field:
	//   {
	//     "response": {this data will be unmarshalled into response arg}
	//   }
	//
	// Nil response arg may be passed to discard response data.
	Exec(method string, params url.Values, response interface{}) error
}

// Options hold configuration data for Client.
type Options struct {
	// AccessToken holds vk.com API access token (optional):
	// https://vk.com/dev/access_token
	AccessToken string
	// Version holds used vk.com API version (strongly recommended):
	// https://vk.com/dev/versions
	Version string
}

// HTTPClient abstracts HTTP client.
type HTTPClient interface {
	PostForm(string, url.Values) (*http.Response, error)
}

// New creates new Client with default HTTP client.
func New(options Options) Client {
	return From(http.DefaultClient, options)
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
	scheme = "https"
	host   = "api.vk.com"
	prefix = "/method/"
)

type responseWrapper struct {
	Error    *Error      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

type client struct {
	options Options
	http    HTTPClient
}

func (c *client) Exec(method string, params url.Values, response interface{}) error {
	if params == nil {
		params = make(url.Values, 2)
	}
	if _, isSet := params["access_token"]; !isSet && c.options.AccessToken != "" {
		params.Set("access_token", c.options.AccessToken)
	}
	if _, isSet := params["v"]; !isSet && c.options.Version != "" {
		params.Set("v", c.options.Version)
	}

	url := urlPool.Get()
	url.Scheme = scheme
	url.Host = host
	url.Path = path.Join(prefix, method)
	defer urlPool.Put(url)

	resp, err := c.http.PostForm(url.String(), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wrapper := responseWrapperPool.Get()
	wrapper.Response = response
	defer responseWrapperPool.Put(wrapper)

	if err := json.NewDecoder(resp.Body).Decode(wrapper); err != nil {
		return err
	}

	if wrapper.Error != nil {
		return *wrapper.Error
	}
	return nil
}
