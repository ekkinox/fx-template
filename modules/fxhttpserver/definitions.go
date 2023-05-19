package fxhttpserver

type MiddlewareDefinition interface {
	Concrete() bool
	Middleware() any
}

type middlewareDefinition struct {
	middleware any
}

func newMiddlewareDefinition(middleware any) *middlewareDefinition {
	return &middlewareDefinition{
		middleware: middleware,
	}
}

func (d *middlewareDefinition) Concrete() bool {
	return isConcreteMiddleware(d.middleware)
}

func (d *middlewareDefinition) Middleware() any {
	return d.middleware
}

type HandlerDefinition interface {
	Concrete() bool
	Method() string
	Path() string
	Handler() any
	Middlewares() []MiddlewareDefinition
}

type handlerDefinition struct {
	method      string
	path        string
	handler     any
	middlewares []MiddlewareDefinition
}

func newHandlerDefinition(method string, path string, handler any, middlewares []MiddlewareDefinition) *handlerDefinition {
	return &handlerDefinition{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

func (d *handlerDefinition) Concrete() bool {
	return isConcreteHandler(d.handler)
}

func (d *handlerDefinition) Method() string {
	return d.method
}

func (d *handlerDefinition) Path() string {
	return d.path
}

func (d *handlerDefinition) Handler() any {
	return d.handler
}

func (d *handlerDefinition) Middlewares() []MiddlewareDefinition {
	return d.middlewares
}
