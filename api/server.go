package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mxrcury/rootgo/router"
	"github.com/mxrcury/rootgo/types"
	"github.com/mxrcury/rootgo/util"
)

type (
	FileType int

	Server struct {
		listenAddr string

		server  *http.Server
		handler *router.Handler

		middlewares []Middleware
		assetsPath  string
	}

	Context struct {
		Request  *http.Request
		Response http.ResponseWriter

		Params *router.Params
		Body   json.Decoder
	}

	Options struct {
		Port string
	}

	Middleware func(http.ResponseWriter, *http.Request)

	IContext interface {
		WriteJSON(interface{}, int)
		WriteError(types.Error)
		Write(interface{}, int) error
		WriteFile(int, []byte, FileType) error
	}
)

func NewServer(r *router.Handler, options Options) *Server {
	listenAddr := fmt.Sprintf(":%s", options.Port)

	return &Server{
		listenAddr: listenAddr,
		server:     &http.Server{Addr: listenAddr, Handler: &router.Handler{Router: r.Router}},
		handler:    r,
	}
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

func (s *Server) USE(middleware func(http.ResponseWriter, *http.Request)) {
	s.middlewares = append(s.middlewares, middleware)
	// Make it working for specific routes
}

func (s *Server) ASSETS(path string) {
	s.handler.ASSETS(path)
}

func (s *Server) Run() error {
	routes := s.handler.Router.Iterate()

	for _, route := range routes {
		fmt.Printf("[%s %s]\n", route.Method, route.Path)
	}

	fmt.Printf("[🔨 Server started on port %s]\n", s.listenAddr[1:])
	return http.ListenAndServe(s.listenAddr, s.server.Handler)
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

const (
	JPEGFileType = iota + 1
	PNGFileType
	SVGFileType
	CSSFileType
	HTMLFileType
	JSFileType
	TTFFileType
	FormDataFileType
	PDFFileType
)

func (c *Context) WriteFile(code int, content []byte, fileType FileType) {
	switch fileType {
	case JPEGFileType:
		c.Response.Header().Set("Content-Type", "image/jpeg")
	case PNGFileType:
		c.Response.Header().Set("Content-Type", "image/png")
	case SVGFileType:
		c.Response.Header().Set("Content-Type", "image/svg+xml")
	case CSSFileType:
		c.Response.Header().Set("Content-Type", "text/css")
	case HTMLFileType:
		c.Response.Header().Set("Content-Type", "text/html")
	case JSFileType:
		c.Response.Header().Set("Content-Type", "application/javascript")
	case PDFFileType:
		c.Response.Header().Set("Content-Type", "application/pdf")
	case TTFFileType:
		c.Response.Header().Set("Content-Type", "font/ttf")
	case FormDataFileType:
		c.Response.Header().Set("Content-Type", "multipart/form-data")
	}
	hashedContent := util.HashValue(content)

	c.Response.Header().Set("Cache-Control", "max-age=3600")

	cacheExpiration := time.Now().Add(time.Hour * 1)
	c.Response.Header().Set("Expires", cacheExpiration.UTC().Format(http.TimeFormat))

	lastModified := time.Now().UTC()
	c.Response.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))

	c.Response.Header().Set("ETag", hashedContent)

	c.Response.WriteHeader(code)

	c.Response.Write(content)
}

func (c *Context) WriteError(err types.Error) {
	c.Response.WriteHeader(err.Status)
	c.Response.Header().Add("Content-Type", "application/json")
	c.Response.Header().Add("Connection", "close")
	json.NewEncoder(c.Response).Encode(err)
}
