package fxhttpserver

type MiddlewareDefinition interface {
	Concrete() bool
	Middleware() any
	Kind() MiddlewareKind
}

type middlewareDefinition struct {
	middleware any
	kind       MiddlewareKind
}

func newMiddlewareDefinition(middleware any, kind MiddlewareKind) *middlewareDefinition {
	return &middlewareDefinition{
		middleware: middleware,
		kind:       kind,
	}
}

func (d *middlewareDefinition) Concrete() bool {
	return isConcreteMiddleware(d.middleware)
}

func (d *middlewareDefinition) Middleware() any {
	return d.middleware
}

func (d *middlewareDefinition) Kind() MiddlewareKind {
	return d.kind
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

type HandlersGroupDefinition interface {
	Prefix() string
	Handlers() []HandlerDefinition
	Middlewares() []MiddlewareDefinition
}

type handlersGroupDefinition struct {
	prefix      string
	handlers    []HandlerDefinition
	middlewares []MiddlewareDefinition
}

func newHandlersGroupDefinition(prefix string, handlers []HandlerDefinition, middlewares []MiddlewareDefinition) *handlersGroupDefinition {
	return &handlersGroupDefinition{
		prefix:      prefix,
		handlers:    handlers,
		middlewares: middlewares,
	}
}

func (h *handlersGroupDefinition) Prefix() string {
	return h.prefix
}

func (h *handlersGroupDefinition) Handlers() []HandlerDefinition {
	return h.handlers
}

func (h *handlersGroupDefinition) Middlewares() []MiddlewareDefinition {
	return h.middlewares
}
