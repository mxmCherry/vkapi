package vkapi

import (
	"net/url"
	"sync"
)

var urlPool = sync.Pool{
	New: func() interface{} {
		return new(url.URL)
	},
}
