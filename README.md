# vkapi [![GoDoc](https://godoc.org/github.com/mxmCherry/vkapi?status.svg)](https://godoc.org/github.com/mxmCherry/vkapi) [![Go Report Card](https://goreportcard.com/badge/github.com/mxmCherry/vkapi)](https://goreportcard.com/report/github.com/mxmCherry/vkapi) [![Build Status](https://travis-ci.org/mxmCherry/vkapi.svg?branch=master)](https://travis-ci.org/mxmCherry/vkapi)

Low-level Go (Golang) vk.com API client


## Example

```go
package main

import "github.com/mxmCherry/vkapi"

// UsersGetRequest represents request params for users.get method: https://vk.com/dev/users.get
type UsersGetRequest struct {
	UserIDs  []uint64 `json:"user_ids"`
	Fields   []string `json:"fields,omitempty"`
	NameCase string   `json:"name_case,omitempty"`

	CaptchaSid string `json:"captcha_sid,omitempty"`
	CaptchaKey string `json:"captcha_key,omitempty"`
}

// UsersGetRequest represents users.get method response: https://vk.com/dev/users.get
type UsersGetResponse []User

// User represents user object: https://vk.com/dev/fields
type User struct {
	ID         uint64 `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ScreenName string `json:"screen_name,omitempty"`
}

func main() {

	// instantiate API client:
	vk := vkapi.New(vkapi.Options{
		AccessToken: "YOUR_ACCESS_TOKEN", // https://vk.com/dev/access_token
	})

	// prepare API request params:
	req := UsersGetRequest{
		UserIDs:  []uint64{1, 2, 3},
		Fields:   []string{"screen_name"},
		NameCase: "nom",
	}

	// prepare container for API response (pointer):
	res := new(UsersGetResponse)

	// execute users.get API method: https://vk.com/dev/users.get
	err := vk.Exec("users.get", vkapi.ToParams(req), res)
	if err != nil {
		// process returned API error: https://vk.com/dev/errors
		if vkErr, ok := err.(vkapi.Error); ok {
			// handle captcha error: https://vk.com/dev/captcha_error
			if vkErr.CaptchaSID != "" {
				panic("Captcha needed: " + vkErr.CaptchaImg)
			}
		}
		// process other (internal) errors:
		panic(err.Error())
	}

	// process response (success):
	for i, user := range *res {
		println(i, user.ID, user.FirstName, user.LastName, user.ScreenName)
	}
}
```
