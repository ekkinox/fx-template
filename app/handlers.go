package app

import (
	"github.com/ekkinox/fx-template/app/handler"
	"github.com/ekkinox/fx-template/app/handler/post"
	"github.com/ekkinox/fx-template/app/middleware"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Options(
		// answer
		fxhttpserver.AsHandler("GET", "/call", handler.NewCallHandler),
		// answer
		fxhttpserver.AsHandler("GET", "/answer", handler.NewAnswerHandler),
		// posts group
		fxhttpserver.AsHandlersGroup(
			"/posts",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", post.NewListPostsHandler),
				fxhttpserver.NewHandlerRegistration("POST", "", post.NewCreatePostHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", post.NewGetPostHandler),
				fxhttpserver.NewHandlerRegistration("PATCH", "/:id", post.NewUpdatePostHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", post.NewDeletePostHandler),
			},
			middleware.NewTestMiddleware,
		),
	)
}
