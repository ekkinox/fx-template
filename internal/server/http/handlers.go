package http

import (
	"github.com/ekkinox/fx-template/internal/server/http/handler/auth"
	crud2 "github.com/ekkinox/fx-template/internal/server/http/handler/crud"
	http2 "github.com/ekkinox/fx-template/internal/server/http/handler/http"
	"github.com/ekkinox/fx-template/internal/server/http/handler/pubsub"
	middleware2 "github.com/ekkinox/fx-template/internal/server/http/middleware"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Options(
		// global
		fxhttpserver.AsMiddleware(middleware2.NewGlobalMiddleware, fxhttpserver.GlobalUse),
		// auth
		fxhttpserver.AsHandler("GET", "/auth/context", auth.NewAuthContextHandler),
		// http
		fxhttpserver.AsHandler("GET", "/http/ping", http2.NewPingHandler, middleware2.NewHandlerMiddleware),
		fxhttpserver.AsHandler("GET", "/http/pong", http2.NewPongHandler),
		// pubsub
		fxhttpserver.AsHandler("GET", "/pubsub/publish", pubsub.NewPublishHandler),
		// crud
		fxhttpserver.AsHandlersGroup(
			"/crud/posts",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", crud2.NewListPostsHandler, middleware2.NewHandlerMiddleware),
				fxhttpserver.NewHandlerRegistration("POST", "", crud2.NewCreatePostHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", crud2.NewGetPostHandler),
				fxhttpserver.NewHandlerRegistration("PATCH", "/:id", crud2.NewUpdatePostHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", crud2.NewDeletePostHandler),
			},
			middleware2.NewGroupMiddleware,
		),
	)
}
