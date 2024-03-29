package sdk

import (
	"net/http"
	"net/url"
)

const (
	// List of available headers
	Authorization string = "Authorization"
	CollectionID  string = "Collection-Id"
)

// Options is a struct containing for customised options for the scrubber API client
type Options struct {
	Headers http.Header
	Query   url.Values
}

// empty Options
func OptInit() *Options {
	return &Options{
		Query:   url.Values{},
		Headers: http.Header{},
	}
}

// Q sets the 'q' Query parameter to the request
func (o *Options) Q(val string) *Options {
	o.Query.Set("q", val)
	return o
}

func setHeaders(req *http.Request, headers http.Header) {
	for name, values := range headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
}
