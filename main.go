package main

import (
	"io"
	"log"
	"net/http"

	"github.com/mxrcury/rootgo/api"
	"github.com/mxrcury/rootgo/config"
	"github.com/mxrcury/rootgo/router"
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
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	r := router.NewRouter(_prefix)
	r.GET("/users/:id", func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		ids := ctx.Params.Get("id")
		io.WriteString(w, "Your user's ID is:"+ids[0])
	})

	r.POST("/users", func(ctx *router.Context, w http.ResponseWriter, r *http.Request) {
		user := new(User)
		err := ctx.Body.Decode(user)

		if err != nil {
			util.WriteError(w, types.Error{Message: "Wrong creation", Status: 400})
		}

		log.Println("USER CREATED:", user)

		io.WriteString(w, "user was successfully created")
	})

	server := api.NewServer(r, api.Options{Port: cfg.Http.Port})

	server.GET("/users/contacts/:id", func(ctx *api.Context) {
		ctx.WriteJSON("ok id:"+ctx.Params.Get("id")[0], 200)
	})

	server.POST("/users/contacts", func(ctx *api.Context) {
		body := new(User)
		ctx.Body.Decode(body)

		ctx.WriteJSON("Successfully", 201)
	})

	server.Run()

	//log.Fatalln(http.ListenAndServe(":"+cfg.Http.Port, r))

	/*

		server := api.NewServer(r, api.Options{Port: cfg.Http.Port})

		r.GET("/users/:id", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, fmt.Sprintf("You get user by ID:%s\n", ctx.Params.Get("id")))
		})

		r.POST("/users/:id", func(c *api.Context, w http.ResponseWriter, r *http.Request) {
			body := new(User)
			err := util.DecodeBody(r.Body).Decode(body)
			if err != nil {
				util.WriteError(w, types.Error{Message: "Bad request", Status: 400})
				return
			}

			util.WriteJSON(w, body, 201)
		})

		r.GET("/users/:id/contacts/:id", func(c *api.Context, w http.ResponseWriter, r *http.Request) {
			params := c.Params.Get("id")
			io.WriteString(w, fmt.Sprintf("You get user by ID:%s, contact id: %s\n", params[0], params[1]))
		})

		r.GET("/users/contacts", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "USER CONTACTS: [ANDRE, PEDRO, LUCAS]")
		})

		log.Printf("Server started on port %s", cfg.Http.Port)

		log.Fatal(server.Run())
	*/
}
