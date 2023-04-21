package handler

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type HelloHandler struct {
	log *zap.Logger
}

func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log: log}
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.log.Info(fmt.Sprintf("called %s", h.Pattern()))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Failed to read request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
		h.log.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (*HelloHandler) Pattern() string {
	return "/hello"
}
