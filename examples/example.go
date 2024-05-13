package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mxrcury/rootgo/pkg/api"
	"github.com/mxrcury/rootgo/pkg/router"
	"github.com/mxrcury/rootgo/types"
)

const _prefix = "/api"

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
  const port = "5000"

	r := router.NewRouter(_prefix)

	server := api.NewServer(r, api.Options{Port: port})

	server.ASSETS("assets")

  r.GET("/users/lol/kek/", func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "success")
  })

	server.GET("/users/:id", func(ctx *api.Context) {
		ids := ctx.Params.Get("id")
		ctx.Write(fmt.Sprintf("<h1>Your user's ID is [%s]</h1><p>You IP is %s</p><b>Magic text: %s</b>\n", ids[0], ctx.Request.RemoteAddr, ctx.Request.URL.Query().Get("msg")), 200)
	})

	server.POST("/users", func(ctx *api.Context) {
		user := new(User)
		err := ctx.Body.Decode(user)

		if err != nil {
			ctx.WriteError(types.Error{Message: "You entered wrong values", Status: 400})
		}

		if err := ctx.Write(user, 200); err != nil {
			ctx.WriteError(types.Error{Message: "Internal server error", Status: 500})
		}
	})

	server.GET("/home", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/index.html")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}
		ctx.WriteFile(200, file, api.HTMLFileType)
	})

	server.GET("/script", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/script.js")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}

		ctx.WriteFile(200, file, api.JSFileType)
	})

	server.Run()
}
