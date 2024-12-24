package sdk

import (
	"github.com/ONSdigital/dis-search-upstream-stub/api"
	"github.com/ONSdigital/dp-net/v2/request"
	"net/http"
	"net/url"
)

const (
	// List of available headers
	Authorization string = request.AuthHeaderKey
)

// Options is a struct containing for customised options for the API client
type Options struct {
	Headers http.Header
	Query   url.Values
}

// Limit sets the 'limit' Query parameter to the request
func (o *Options) Limit(val string) *Options {
	if o.Query == nil {
		var queryParams url.Values = make(map[string][]string)
		o.Query = queryParams
	}
	_, paramExists := o.Query[api.ParamLimit]
	if paramExists {
		o.Query.Set(api.ParamLimit, val)
	} else {
		o.Query.Add(api.ParamLimit, val)
	}
	return o
}

// Offset sets the 'offset' Query parameter to the request
func (o *Options) Offset(val string) *Options {
	if o.Query == nil {
		var queryParams url.Values = make(map[string][]string)
		o.Query = queryParams
	}
	_, paramExists := o.Query[api.ParamOffset]
	if paramExists {
		o.Query.Set(api.ParamOffset, val)
	} else {
		o.Query.Add(api.ParamOffset, val)
	}
	return o
}

func setHeaders(req *http.Request, headers http.Header) {
	for name, values := range headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
}
