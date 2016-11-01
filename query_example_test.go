package vkapi_test

import (
	"fmt"

	"github.com/mxmCherry/vkapi"
)

func ExampleQuery() {
	type UsersGetRequest struct {
		UserIDs  []uint64 `json:"user_ids,omitempty"`
		Fields   []string `json:"fields,omitempty"`
		NameCase string   `json:"name_case,omitempty"`
	}

	request := UsersGetRequest{
		UserIDs:  []uint64{111, 222, 333},
		Fields:   []string{"first_name", "last_name", "screen_name"},
		NameCase: "nom",
	}

	query := vkapi.Query(request)

	fmt.Println(query.Encode())
	// Output: fields=first_name%2Clast_name%2Cscreen_name&name_case=nom&user_ids=111%2C222%2C333
}
