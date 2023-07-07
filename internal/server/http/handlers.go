package http

import (
	"github.com/ekkinox/fx-template/internal/server/http/handler/auth"
	"github.com/ekkinox/fx-template/internal/server/http/handler/http"
	"github.com/ekkinox/fx-template/internal/server/http/handler/posts"
	"github.com/ekkinox/fx-template/internal/server/http/handler/pubsub"
	"github.com/ekkinox/fx-template/internal/server/http/middleware"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Options(
		// global
		fxhttpserver.AsMiddleware(middleware.NewGlobalMiddleware, fxhttpserver.GlobalUse),
		// auth
		fxhttpserver.AsHandler("GET", "/auth", auth.NewAuthContextHandler),
		// http
		fxhttpserver.AsHandler("GET", "/ping", http.NewPingHandler, middleware.NewHandlerMiddleware),
		fxhttpserver.AsHandler("GET", "/pong", http.NewPongHandler),
		fxhttpserver.AsHandler("GET", "/test", http.NewTestHandler),
		// pubsub
		fxhttpserver.AsHandler("GET", "/publish", pubsub.NewPublishHandler),
		// posts crud
		fxhttpserver.AsHandlersGroup(
			"/posts",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", posts.NewListPostsHandler, middleware.NewHandlerMiddleware),
				fxhttpserver.NewHandlerRegistration("POST", "", posts.NewCreatePostHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", posts.NewGetPostHandler),
				fxhttpserver.NewHandlerRegistration("PATCH", "/:id", posts.NewUpdatePostHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", posts.NewDeletePostHandler),
			},
			middleware.NewGroupMiddleware,
		),
	)
}
