package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mxrcury/rootgo/api"
	"github.com/mxrcury/rootgo/config"
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

	r.GET("/users/:id/contacts/:email", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("You get user by ID:%s, contact id: %s\n", ctx.Params.Get("id"), ctx.Params.Get("email")))
	})

	r.GET("/users/contacts", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "USER CONTACTS: [ANDRE, PEDRO, LUCAS]")
	})

	log.Printf("Server started on port %s", cfg.Http.Port)

	log.Fatal(server.Run())
}
