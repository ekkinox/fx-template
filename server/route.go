package server

import (
	"net/http"

	"go.uber.org/fx"
)

type Route interface {
	http.Handler
	Pattern() string
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
