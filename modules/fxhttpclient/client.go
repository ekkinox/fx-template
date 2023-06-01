package fxhttpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const HttpClientContextRequestKey = "_httpClientContextRequest"

type HttpClient struct {
	context          context.Context
	client           *http.Client
	headersToForward map[string][]string
}

func NewCtxHttpClient(ctx context.Context, opts ...HttpClientOption) *HttpClient {

	appliedOpts := defaultHttpClientOptions
	for _, applyOpt := range opts {
		applyOpt(&appliedOpts)
	}

	client := &http.Client{
		Transport:     otelhttp.NewTransport(appliedOpts.Transport),
		CheckRedirect: appliedOpts.CheckRedirect,
		Timeout:       appliedOpts.Timeout,
	}

	return &HttpClient{
		context:          ctx,
		client:           client,
		headersToForward: appliedOpts.HeadersToForward,
	}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	req = req.WithContext(c.context)

	fmt.Printf("in DO !\n**********************************\n")
	fmt.Printf("ctx: %+v", c.context)

	for name, values := range c.headersToForward {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	return c.client.Do(req)
}

func (c *HttpClient) Get(url string) (resp *http.Response, err error) {

	fmt.Printf("in GET !\n**********************************\n")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HttpClient) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

func (c *HttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
