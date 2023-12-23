package api

import (
	"io"
	"net/http"
	"slices"
	"strings"
)

type IHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Handler struct {
	Body interface{}

	Router *Router
}

type (
	Router struct {
		node *Node
	}

	Node struct {
		Path string

		Handlers map[string]func(*Context, http.ResponseWriter, *http.Request)
		Children map[string]*Node
	}

	Context struct {
		Params *Params
	}

	Params struct {
		Values map[string]string
	}
)

func NewRouter(path string) *Router {
	return &Router{
		node: &Node{
			Path:     path,
			Handlers: make(map[string]func(*Context, http.ResponseWriter, *http.Request)),
			Children: make(map[string]*Node),
		},
	}
}

func NewParams() *Params {
	return &Params{make(map[string]string)}
}

func newContext(params *Params) *Context {
	return &Context{Params: params}
}

/*

{
	Path: "/api"
	Handlers {}
	Children [
		"/users" {
			Path: "/users",
			Handlers {
				"GET": func
				"POST": func
			}
			Children {
				"/:id": {
					Path: "/:id"
					Handlers {
						"GET": func
					}
					Children {}
				}
			}
		}
	]
}
*/

func (r *Router) add(method, path string, handler func(*Context, http.ResponseWriter, *http.Request)) {
	type methods []string

	allowedMethods := methods{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodHead, http.MethodPut}
	if !slices.Contains[methods, string](allowedMethods, method) {
		return
	}
	explodedPath := explodePath(path)
	node := r.node
	for index, path := range explodedPath {
		isLastElement := index == len(explodedPath)-1
		if node.Children == nil {
			node.Children = make(map[string]*Node)
		}
		if node.Handlers == nil {
			node.Handlers = make(map[string]func(*Context, http.ResponseWriter, *http.Request))
		}
		if node.Children[path] == nil {
			node.Children[path] = &Node{
				Path: path,
			}
			node = node.Children[path]
		} else {
			node = node.Children[path]
			if !isLastElement {
				continue
			}
		}
		if isLastElement {
			if node.Children == nil {
				node.Children = make(map[string]*Node)
			}
			if node.Handlers == nil {
				node.Handlers = map[string]func(*Context, http.ResponseWriter, *http.Request){
					method: handler,
				}
			} else {
				node.Handlers[method] = handler
			}
		}

	}
}

func (r *Router) search(method, path string) (func(*Context, http.ResponseWriter, *http.Request), *Params) {
	params := NewParams()
	if !strings.HasPrefix(path, r.node.Path) {
		return nil, nil
	}

	explodedPath := explodePath(strings.Replace(path, r.node.Path, "", 1))
	node := r.node
	for index, path := range explodedPath {
		isLastElement := index == len(explodedPath)-1
		if node.Children[path] != nil && !isLastElement {
			node = node.Children[path]
			continue
		}
		if node.Children[path] != nil && isLastElement {
			if node.Children[path].Handlers[method] == nil {
				return nil, nil
			}
			return node.Children[path].Handlers[method], params
		}

		if node.Children[path] == nil && !isLastElement {
			for _, v := range node.Children {
				if strings.HasPrefix(v.Path, ":") {
					node = node.Children[v.Path]
					params.Set(strings.Replace(v.Path, ":", "", 1), path)
					continue
				}
			}
		}

		if node.Children[path] == nil && isLastElement {
			for _, v := range node.Children {
				if strings.HasPrefix(v.Path, ":") {
					params.Set(strings.Replace(v.Path, ":", "", 1), path)
					return node.Children[v.Path].Handlers[method], params
				}
			}
		}
	}

	return nil, nil
}

func explodePath(path string) []string {
	explodedPath := new([]string)
	splittedPath := strings.Split(path, "/")

	for _, val := range splittedPath {
		if val != "" {
			*explodedPath = append(*explodedPath, val)
		}
	}
	return *explodedPath
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params := h.Router.search(r.Method, r.URL.Path)
	if handler != nil {
		ctx := newContext(params)
		handler(ctx, w, r)
	} else {
		io.WriteString(w, "404 not found page")
	}
}
func (r *Router) GET(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	r.add("GET", path, handler)
}

func (r *Router) POST(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	r.add("POST", path, handler)
}

func (r *Router) PUT(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	r.add("PUT", path, handler)
}

func (r *Router) PATCH(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	r.add("PATCH", path, handler)
}

func (r *Router) DELETE(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	r.add("DELETE", path, handler)
}

func (p *Params) Set(key, value string) {
	p.Values[key] = value
}

func (p *Params) Get(key string) string {
	if _, ok := p.Values[key]; ok {
		return p.Values[key]
	}
	return ""
}
