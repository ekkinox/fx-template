package fxhttpclient

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpClient struct {
	*http.Client
	context context.Context
	headers []string
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
		client,
		ctx,
		appliedOpts.Headers,
	}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	req = req.WithContext(c.context)

	c.context.(echo.Context).Logger().Info("in do !\n**********************************\n")
	fmt.Printf("in do !\n**********************************\n")

	for _, h := range c.headers {
		req.Header.Add(h, req.Header.Get(h))
		fmt.Printf("adding header %s: %s", h, req.Header.Get(h))
	}

	return c.Client.Do(req.WithContext(c.context))
}
