package fxhttpclient

import (
	"net/http"
	"strings"
	"time"
)

var DefaultHeadersToForward = []string{
	"Authorization",
	"X-Request-ID",
	"aaa",
}

type options struct {
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	Timeout          time.Duration
	HeadersToForward map[string][]string
}

var defaultHttpClientOptions = options{
	Transport:        http.DefaultTransport,
	CheckRedirect:    nil,
	Timeout:          time.Second * 10,
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

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if strings.ToLower(v) == strings.ToLower(s) {
			return true
		}
	}

	return false
}
