package main

import (
	"io"
	"log"
	"net/http"

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

	log.Fatalln(http.ListenAndServe(":"+cfg.Http.Port, r))
}
