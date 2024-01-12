package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mxrcury/rootgo/api"
	"github.com/mxrcury/rootgo/config"
	"github.com/mxrcury/rootgo/router"
	"github.com/mxrcury/rootgo/types"
)

const _prefix = "/api"

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	r := router.NewRouter(_prefix)

	server := api.NewServer(r, api.Options{Port: cfg.Http.Port})

	server.ASSETS("assets")

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

	server.GET("/ok", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/index.html")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}
		ctx.WriteFile(200, file, api.HTMLFileType)
	})

	server.GET("/script.js", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/script.js")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}

		ctx.WriteFile(200, file, api.JSFileType)
	})

	server.GET("/document", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/doc.pdf")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}

		ctx.WriteFile(200, file, api.PDFFileType)
	})

	server.Run()
}
