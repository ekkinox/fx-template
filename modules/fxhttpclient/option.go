package fxhttpclient

import (
	"net/http"
	"strings"
	"time"
)

var DefaultHeadersToForward = []string{
	"authorization",
	"x-request-id",
	"traceparent",
}

type options struct {
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	Jar              http.CookieJar
	Timeout          time.Duration
	HeadersToForward map[string][]string
}

var defaultHttpClientOptions = options{
	Transport:        boostedTransport(),
	CheckRedirect:    nil,
	Timeout:          time.Second * 10,
	Jar:              nil,
	HeadersToForward: map[string][]string{},
}

type HttpClientOption func(o *options)

func WithTransport(t http.RoundTripper) HttpClientOption {
	return func(o *options) {
		o.Transport = t
	}
}

func WithCheckRedirect(f func(req *http.Request, via []*http.Request) error) HttpClientOption {
	return func(o *options) {
		o.CheckRedirect = f
	}
}

func WithCookieJar(j http.CookieJar) HttpClientOption {
	return func(o *options) {
		o.Jar = j
	}
}

func WithTimeout(t time.Duration) HttpClientOption {
	return func(o *options) {
		o.Timeout = t
	}
}

func WithRequestHeadersToForward(req *http.Request, headersNames []string) HttpClientOption {
	headersToForward := map[string][]string{}

	for name, values := range req.Header {
		if contains(headersNames, name) {
			headersToForward[strings.ToLower(name)] = values
		}
	}

	return func(o *options) {
		o.HeadersToForward = headersToForward
	}
}

func boostedTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()

	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return t
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if strings.ToLower(v) == strings.ToLower(s) {
			return true
		}
	}

	return false
}
