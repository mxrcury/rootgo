package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mxrcury/rootgo/router"
	"github.com/mxrcury/rootgo/types"
)

type Server struct {
	listenAddr string

	server  *http.Server
	handler *router.Handler
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter

	Params *router.Params
	Body   json.Decoder
}

type IContext interface {
	WriteJSON(interface{}, int)
	WriteError(types.Error)
}

type Options struct {
	Port string
}

func NewServer(r *router.Handler, options Options) *Server {
	listenAddr := fmt.Sprintf(":%s", options.Port)

	return &Server{
		listenAddr: listenAddr,
		server:     &http.Server{Addr: listenAddr, Handler: &router.Handler{Router: r.Router}},
		handler:    r,
	}
}

func (c *Context) Write(data interface{}, status int) {
	switch v := data.(type) {
	case string:
		if isValidJSON := json.Valid([]byte(v)); isValidJSON { // change check of every incoming type if it's valid json
			c.WriteJSON(data, status)
			return
		} else {
			c.Response.WriteHeader(status)
			c.Response.Header().Add("Content-Type", "text/plain")
			c.Response.Header().Add("Connection", "close")
			io.WriteString(c.Response, v)
		}
	}
}

func (c *Context) WriteJSON(data interface{}, status int) {
	c.Response.WriteHeader(status)
	c.Response.Header().Add("Content-Type", "application/json")
	c.Response.Header().Add("Connection", "close")
	json.NewEncoder(c.Response).Encode(data)
}

func (c *Context) WriteError(err types.Error) {
	c.Response.WriteHeader(err.Status)
	c.Response.Header().Add("Content-Type", "application/json")
	c.Response.Header().Add("Connection", "close")
	json.NewEncoder(c.Response).Encode(err)
}

func (s *Server) GET(path string, handler func(*Context)) {
	s.handler.Router.Add("GET", path, func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		handler(&Context{Request: r, Response: w, Params: ctx.Params, Body: *ctx.Body})
	})
}

func (s *Server) POST(path string, handler func(*Context)) {
	s.handler.Router.Add("POST", path, func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		handler(&Context{Request: r, Response: w, Params: ctx.Params, Body: *ctx.Body})
	})
}

func (s *Server) PUT(path string, handler func(*Context)) {
	s.handler.Router.Add("PUT", path, func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		handler(&Context{Request: r, Response: w, Params: ctx.Params, Body: *ctx.Body})
	})
}

func (s *Server) PATCH(path string, handler func(*Context)) {
	s.handler.Router.Add("PATCH", path, func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		handler(&Context{Request: r, Response: w, Params: ctx.Params, Body: *ctx.Body})
	})
}

func (s *Server) DELETE(path string, handler func(*Context)) {
	s.handler.Router.Add("DELETE", path, func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		handler(&Context{Request: r, Response: w, Params: ctx.Params, Body: *ctx.Body})
	})
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddr, s.server.Handler)
}
