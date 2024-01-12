package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/mxrcury/rootgo/util"
)

type IHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	GET(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request))
	POST(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request))
	PATCH(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request))
	PUT(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request))
	DELETE(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request))
}

type (
	Handler struct {
		Router *Router
	}

	Router struct {
		node *Node

		assets *Assets
	}

	Node struct {
		Path     string
		FullPath string

		Handlers map[string]func(*Context, http.ResponseWriter, *http.Request)
		Children map[string]*Node
	}

	Context struct {
		Params *Params

		Body *json.Decoder
	}

	Params struct {
		Values map[string][]string
	}

	Body struct {
		BodyDecoder *json.Decoder
	}

	Assets struct {
		Path string
	}
)

func NewRouter(prefix string) *Handler {
	return &Handler{
		Router: &Router{
			node: &Node{
				Path:     prefix,
				FullPath: prefix,
				Handlers: make(map[string]func(*Context, http.ResponseWriter, *http.Request)),
				Children: make(map[string]*Node),
			},
		},
	}
}

func newParams() *Params {
	return &Params{make(map[string][]string)}
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

func (r *Router) Add(method, path string, handler func(*Context, http.ResponseWriter, *http.Request)) {
	type methods []string

	allowedMethods := methods{http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodHead, http.MethodPut}
	if !slices.Contains[methods, string](allowedMethods, method) {
		return
	}
	explodedPath := explodePath(path)
	node := r.node
	prefix := node.Path
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
				Path:     path,
				FullPath: fmt.Sprintf("%s/%s", prefix, strings.Join(explodedPath, "/")),
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
	params := newParams()
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
					params.set(strings.Replace(v.Path, ":", "", 1), path)
					continue
				}
			}
		}

		if node.Children[path] == nil && isLastElement {
			for _, v := range node.Children {
				if strings.HasPrefix(v.Path, ":") {
					params.set(strings.Replace(v.Path, ":", "", 1), path)
					return node.Children[v.Path].Handlers[method], params
				}
			}
		}
	}

	return nil, nil
}

type Route struct {
	Method string
	Path   string
}

func iterate(list *[]*Route, node *Node) {
	if len(node.Children) != 0 {
		if len(node.Handlers) != 0 {
			for h := range node.Handlers {
				route := &Route{Method: h, Path: node.FullPath}
				*list = append(*list, route)
			}
		}
		for c := range node.Children {
			iterate(list, node.Children[c])
		}

	} else {
		for h := range node.Handlers {
			route := &Route{Method: h, Path: node.FullPath}
			*list = append(*list, route)
		}
	}
	return
}

func (r *Router) Iterate() []*Route {
	routes := []*Route{}
	node := r.node

	iterate(&routes, node)

	return routes
}

func explodePath(path string) []string {
	explodedPath := make([]string, 0, 6)
	splitPath := strings.Split(path, "/")

	for _, val := range splitPath {
		if val != "" {
			explodedPath = append(explodedPath, val)
		}
	}
	return explodedPath
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params := h.Router.search(r.Method, r.URL.Path)

	if handler != nil {
		ctx := newContext(params)
		ctx.NewBodyDecoder(r.Body)
		handler(ctx, w, r)
		return
	}

	areAssetsEnabled := h.Router.assets != nil && strings.TrimSpace(h.Router.assets.Path) != ""
	if areAssetsEnabled && strings.HasPrefix(r.URL.Path, h.Router.assets.Path) || strings.HasPrefix(strings.ReplaceAll(r.URL.Path, "/", ""), strings.ReplaceAll(h.Router.assets.Path, "/", "")) {

		assetsPath := path.Join(strings.Split(r.URL.Path, "/")...)

		file, err := os.ReadFile(assetsPath)
		if err != nil {
			json.NewEncoder(w).Encode(util.NewError("file not found", 404))
			return
		}

		w.Write(file)
		return
	}
	json.NewEncoder(w).Encode(util.NewError("not found", 404))
}

func (h *Handler) GET(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	h.Router.Add("GET", path, handler)
}

func (h *Handler) POST(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	h.Router.Add("POST", path, handler)
}

func (h *Handler) PUT(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	h.Router.Add("PUT", path, handler)
}

func (h *Handler) PATCH(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	h.Router.Add("PATCH", path, handler)
}

func (h *Handler) DELETE(path string, handler func(ctx *Context, w http.ResponseWriter, r *http.Request)) {
	h.Router.Add("DELETE", path, handler)
}

func (h *Handler) ASSETS(path string) {
	h.Router.assets = &Assets{Path: path}
}

func (p *Params) set(key, value string) {
	if _, ok := p.Values[key]; ok {
		p.Values[key] = append(p.Values[key], value)
	} else {
		p.Values[key] = []string{value}
	}
}

func (p *Params) Get(key string) []string {
	if _, ok := p.Values[key]; ok {
		return p.Values[key]
	}
	return []string{}
}

func (c *Context) NewBodyDecoder(reader io.ReadCloser) {
	c.Body = json.NewDecoder(reader)
}
