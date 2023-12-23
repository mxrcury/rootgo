package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	listenAddr string

	server *http.Server
}

func NewServer(listenAddr string) *Server {
	listenAddr = fmt.Sprintf(":%s", listenAddr)

	return &Server{
		listenAddr: listenAddr,
		server:     &http.Server{Addr: listenAddr},
	}
}

func (s *Server) Router(router *Router) {
	s.server.Handler = &Handler{Router: router}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddr, s.server.Handler)
}
