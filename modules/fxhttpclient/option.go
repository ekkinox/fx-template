package fxhttpclient

import (
	"net/http"
	"time"
)

var DefaultHeadersToForward = []string{"authorization"}

type options struct {
	Transport     http.RoundTripper
	CheckRedirect func(req *http.Request, via []*http.Request) error
	Timeout       time.Duration
	Headers       []string
}

var defaultHttpClientOptions = options{
	Transport:     http.DefaultTransport,
	CheckRedirect: nil,
	Timeout:       time.Second * 10,
	Headers:       DefaultHeadersToForward,
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

func WithHeaders(h []string) HttpClientOption {
	return func(o *options) {
		o.Headers = h
	}
}
