package ws

import (
	"fmt"
	"net/http"
)

type httpHandler struct {
	handler http.HandlerFunc
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(w, r)
}

type startable interface {
	Start() http.HandlerFunc
}

func httpServer(port int, ws startable) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: &httpHandler{ws.Start()},
	}

	go server.ListenAndServe()

	return server
}
