package api

import (
	"context"
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
		Node *Node
	}

	Node struct {
		Path string

		Handlers map[string]func(context.Context, http.ResponseWriter, *http.Request)
		Children map[string]*Node
	}
)

func NewRouter(path string) *Router {
	return &Router{
		Node: &Node{
			Path:     path,
			Handlers: make(map[string]func(context.Context, http.ResponseWriter, *http.Request)),
			Children: make(map[string]*Node),
		},
	}
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

func (r *Router) add(method, path string, handler func(context.Context, http.ResponseWriter, *http.Request)) {
	type methods []string

	allowedMethods := methods{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodHead, http.MethodPut}
	if !slices.Contains[methods, string](allowedMethods, method) {
		return
	}
	explodedPath := explodePath(path)
	node := r.Node
	for index, path := range explodedPath {
		isLastElement := index == len(explodedPath)-1
		if node.Children == nil {
			node.Children = make(map[string]*Node)
		}
		if node.Handlers == nil {
			node.Handlers = make(map[string]func(context.Context, http.ResponseWriter, *http.Request))
		}
		if node.Children[path] == nil {
			node.Children[path] = &Node{
				Path: path,
			}
		} else {
			node = node.Children[path]
			if !isLastElement {
				continue
			}
		}
		if isLastElement {
			if path != node.Path {
				node = node.Children[path]
			}
			if node.Handlers == nil {
				node.Handlers = map[string]func(context.Context, http.ResponseWriter, *http.Request){
					method: handler,
				}
			} else {
				node.Handlers[method] = handler
			}
		}

	}
}

func (r *Router) search(method, path string) func(context.Context, http.ResponseWriter, *http.Request) {
	explodedPath := explodePath(strings.Replace(path, "/api", "", 1))
	node := r.Node
	for index, path := range explodedPath {
		isLastElement := index == len(explodedPath)-1
		if node.Children[path] != nil && !isLastElement {
			if node.Children[path] == nil {
				return nil
			}
			node = node.Children[path]
		}

		if node.Children[path] != nil && isLastElement {
			if node.Children[path].Handlers[method] == nil {
				return nil
			}
			return node.Children[path].Handlers[method]
		}
	}

	return nil
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
	handler := h.Router.search(r.Method, r.URL.Path)
	if handler != nil {
		handler(context.Background(), w, r)
	} else {
		io.WriteString(w, "404 not found page")
	}
}
func (r *Router) GET(path string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	r.add("GET", path, handler)
}

func (r *Router) POST(path string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	r.add("POST", path, handler)
}

func (r *Router) PUT(path string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	r.add("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	r.add("DELETE", path, handler)
}
