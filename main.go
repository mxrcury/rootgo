package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mxrcury/rootgo/api"
	"github.com/mxrcury/rootgo/config"
	"github.com/mxrcury/rootgo/types"
	"github.com/mxrcury/rootgo/util"
)

const _prefix = "/api"

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {

	r := api.NewRouter("/api")

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	server := api.NewServer(cfg.Http.Port)

	server.Router(r)

	r.GET("/users/:id", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("You get user by ID:%s\n", ctx.Params.Get("id")))
	})

	r.POST("/users/:id", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		body := new(User)
		err := util.DecodeBody(r.Body).Decode(body)
		if err != nil {
			util.WriteError(w, types.Error{Message: "Bad request", Status: 400})
			return
		}

		util.WriteJSON(w, body, 201)
	})

	r.GET("/users/:id/contacts/:email", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("You get user by ID:%s, contact id: %s\n", ctx.Params.Get("id"), ctx.Params.Get("email")))
	})

	r.GET("/users/contacts", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "USER CONTACTS: [ANDRE, PEDRO, LUCAS]")
	})

	log.Printf("Server started on port %s", cfg.Http.Port)

	log.Fatal(server.Run())
}
