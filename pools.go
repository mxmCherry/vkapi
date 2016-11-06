package vkapi

import (
	"net/url"
	"sync"
)

var (
	urlPool             = newURLPooler()
	responseWrapperPool = newResponseWrapperPooler()
)

// ----------------------------------------------------------------------------

type urlPooler struct {
	pool sync.Pool
}

func newURLPooler() urlPooler {
	return urlPooler{
		pool: sync.Pool{
			New: func() interface{} {
				return new(url.URL)
			},
		},
	}
}

func (p *urlPooler) Get() *url.URL {
	return p.pool.Get().(*url.URL)
}

func (p *urlPooler) Put(u *url.URL) {
	*u = url.URL{}
	p.pool.Put(u)
}

// ----------------------------------------------------------------------------

type responseWrapperPooler struct {
	pool sync.Pool
}

func newResponseWrapperPooler() responseWrapperPooler {
	return responseWrapperPooler{
		pool: sync.Pool{
			New: func() interface{} {
				return new(responseWrapper)
			},
		},
	}
}

func (p *responseWrapperPooler) Get() *responseWrapper {
	return p.pool.Get().(*responseWrapper)
}

func (p *responseWrapperPooler) Put(w *responseWrapper) {
	*w = responseWrapper{}
	p.pool.Put(w)
}
