# Root
Simple HTTP router/server for Golang

### Current features overview
- Adding `GET`, `POST`, `PUT`, `DELETE` endpoints
- Routes with params
- Getting value of params
- Separate server that the router can be connected to

### Coming features
- Make router fully separate from server and possible of using without need to use server
- Add possibility to have values of several routes with similar values
- Connect server utils like `WriteJSON`, `DecodeBody` as methodss of the server(or context?)


### Usage example
```go

func main() {

	r := api.NewRouter("/api")

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	server := api.NewServer(cfg.Http.Port)

	server.Router(r)

	/*
	GET /users/:id
	curl -X GET http://localhost:8000/api/users/31231232
	---
	'You get user by ID:31231232'
	*/
	r.GET("/users/:id", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("You get user by ID:%s\n", ctx.Params.Get("id")))
	})

	/*
	POST /users/:id
	curl -X POST -H "Content-Type: application/json" -d '{"name": "Anthony", "age": 19, "city": "Miami"}' http://localhost:8000/api/users/dwqdqw
	---
	'{"name":"Anthony","age":19,"city":"Miami"}'
	*/
	r.POST("/users/:id", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		body := new(User)
		err := util.DecodeBody(r.Body).Decode(body)
		if err != nil {
			util.WriteError(w, types.Error{Message: "Bad request", Status: 400})
			return
		}

		util.WriteJSON(w, body, 201)
	})

	/*
	GET /users/:id/contacts/:email
	curl -X GET http://localhost:8000/api/users/31231232/contacts/test@gmail.com
	---
	'You get user by ID:31231232, contact id: test@gmail.com'
	*/
	r.GET("/users/:id/contacts/:email", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("You get user by ID:%s, email: %s\n", ctx.Params.Get("id"), ctx.Params.Get("email")))
	})

	/*
	GET /users/contacts
	curl -X GET http://localhost:8000/api/users/contacts
	---
	'USER CONTACTS: [ANDRE, PEDRO, LUCAS]'
	*/
	r.GET("/users/contacts", func(ctx *api.Context, w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "USER CONTACTS: [ANDRE, PEDRO, LUCAS]")
	})

	log.Printf("Server started on port %s", cfg.Http.Port)

	log.Fatal(server.Run())
}
```

Licensed by MIT
