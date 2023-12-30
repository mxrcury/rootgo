package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *Context) Write(data interface{}, status int) error {
	switch d := data.(type) {
	case string:
		c.Response.WriteHeader(status)
		c.Response.Header().Add("Content-Type", "text/html")
		c.Response.Header().Add("Connection", "close")
		_, err := io.WriteString(c.Response, d)
		return err
	default:
		return c.WriteJSON(data, status)
	}
}

func (c *Context) WriteJSON(data interface{}, status int) error {
	c.Response.WriteHeader(status)
	c.Response.Header().Add("Content-Type", "application/json")
	c.Response.Header().Add("Connection", "close")
	err := json.NewEncoder(c.Response).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

type FileType int

const (
	JPEGType = iota + 1
	PNGType
	SVGType
	CSSType
	JSType
	TTFType
	FormDataType
	PDFType
)

func (c *Context) WriteFile(content []byte, fileType FileType) {

	file := new([]byte)
	buff := make([]byte, 512)
	f := bytes.NewReader(content)
	for {
		n, err := f.Read(buff)
		if err != nil {
			break
		}
		_ = n
		*file = append(*file, buff[:n]...)
	}
	switch fileType {
	case JPEGType:

		c.Response.Header().Add("Content-Type", "image/jpeg")
		c.Response.Write(*file)
	case PNGType:
		c.Response.Header().Add("Content-Type", "image/png")
	case JSType:
		c.Response.Header().Add("Content-Type", "application/javascript")
		log.Println("JS TYPE")
	}
	c.Response.Header().Add("Content-Length", fmt.Sprintf("%d", len(*file)))
	c.Response.Write(*file)
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
	// (?): Maybe log it only when logger is enabled
	log.Printf("🔨 Server started on port %s\n", s.listenAddr[1:])
	return http.ListenAndServe(s.listenAddr, s.server.Handler)
}
