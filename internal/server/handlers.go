package server

import (
	"github.com/ekkinox/fx-template/internal/server/handler/crud"
	"github.com/ekkinox/fx-template/internal/server/handler/http"
	"github.com/ekkinox/fx-template/internal/server/handler/pubsub"
	"github.com/ekkinox/fx-template/internal/server/middleware"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Options(
		// http
		fxhttpserver.AsHandler("GET", "/http/ping", http.NewPingHandler),
		fxhttpserver.AsHandler("GET", "/http/pong", http.NewPongHandler),
		// pubsub
		fxhttpserver.AsHandler("GET", "/pubsub/publish", pubsub.NewPublishHandler),
		// crud
		fxhttpserver.AsHandlersGroup(
			"/crud/posts",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", crud.NewListPostsHandler),
				fxhttpserver.NewHandlerRegistration("POST", "", crud.NewCreatePostHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", crud.NewGetPostHandler),
				fxhttpserver.NewHandlerRegistration("PATCH", "/:id", crud.NewUpdatePostHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", crud.NewDeletePostHandler),
			},
			middleware.NewTestMiddleware,
		),
	)
}
