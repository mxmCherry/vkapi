package vkapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type Options struct {
	AccessToken string
}

type Client interface {
	Exec(method string, query url.Values, response interface{}) error
}

func New(options Options) Client {
	return &client{
		options: options,
		http:    new(http.Client),
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
	http    interface {
		Get(string) (*http.Response, error)
	}
}

func (c *client) Exec(method string, query url.Values, response interface{}) error {
	if query == nil {
		query = url.Values{}
	}
	if c.options.AccessToken != "" {
		query.Set("access_token", c.options.AccessToken)
	}
	query.Set("v", version)

	resp, err := c.http.Get((&url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path.Join(prefix, method),
		RawQuery: query.Encode(),
	}).String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wrapper := &struct {
		Error *struct {
			ErrorCode uint64 `json:"error_code"`
			ErrorMsg  string `json:"error_msg"`
		} `json:"error,omitempty"`
		Response interface{} `json:"response"`
	}{
		Response: response,
	}
	if err := json.NewDecoder(resp.Body).Decode(wrapper); err != nil {
		return err
	}

	if wrapper.Error != nil {
		return fmt.Errorf("vkapi: %s (code %d)", wrapper.Error.ErrorMsg, wrapper.Error.ErrorCode)
	}
	return nil
}
