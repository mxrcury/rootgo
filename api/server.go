package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	listenAddr string

	server *http.Server
}

type Options struct {
	Port string
}

func NewServer(handler http.Handler, options Options) *Server {
	listenAddr := fmt.Sprintf(":%s", options.Port)

	return &Server{
		listenAddr: listenAddr,
		server:     &http.Server{Addr: listenAddr, Handler: handler},
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddr, s.server.Handler)
}
