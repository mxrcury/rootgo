# Root
Simple and lightweight HTTP router/server for Golang
```
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣷⠀⠀⠀⢀⣤⣾⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣀⣀⣀⠀⠀⠀⠹⣧⡀⠀⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⣷⣦⣄⠘⣷⡀⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠛⠛⠀⠘⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣹⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡿⠀⢴⣾⠿⠛⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⢲⣶⣶⣶⡖⠒⠒⠒⠒⠒⠒⠒⠒⠒⠒⠀⠀
⠀⠀⠀⠀⠀⢠⡶⠾⠆⠀⠀⠀⠀⣸⣿⠁⠈⠙⢷⡶⠶⣦⣄⡀⠀⠀⠀⠀⠀⠀
⠀⠸⠷⣦⡀⠘⢷⣄⠀⠀⠀⠀⢠⡿⣯⡀⠀⠀⠈⢿⡄⠀⠉⠁⠀⠀⣠⡄⠀⠀
⠀⠀⢠⡾⠿⠛⠻⠿⠿⠶⠶⠾⠛⢁⣸⣷⣦⣀⠀⠈⠁⠀⠀⢀⣤⣾⡋⠀⠀⠀
⠀⠀⣿⠁⠀⠀⠀⠀⠀⠀⠀⢀⣴⠿⣯⣀⡀⠉⠙⠛⠓⠒⠚⠋⠁⠈⢿⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠀
```
### Current features overview
- Adding `GET`, `POST`, `PATCH`, `PUT`, `DELETE` endpoints
- Routes with params
- Getting value of params, several params with the same name presented as slice
- Router can be used both with standard library and our server wrapper
- Server wrapper that implements some high level utils(Write, WriteFiles)
- Working with views

### Coming features
- [Server] Implement caching(in memory and in redis or any other db that implement Cache interface in `/caching/cache.go`) 

### Usage example

```go
func main(){
	r := router.NewRouter("/api")

	server := api.NewServer(r, api.Options{Port: "8000"})

	server.ASSETS("/assets")

	server.GET("/users/contacts/:id", func(ctx *api.Context) {
		ids := ctx.Params.Get("id")
		ctx.WriteJSON("GET id:"+id[0], 200)
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

	server.GET("/document", func(ctx *api.Context) {
		file, err := os.ReadFile("./assets/doc.pdf")
		if err != nil {
			ctx.WriteError(types.Error{Message: err.Error(), Status: 500})
		}

		ctx.WriteFile(200, file, api.PDFFileType)
	})

	server.Run()
}

```

<!--
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
-->

### Usage only router without server wrapper

```golang
func main(){
	r := router.NewRouter("/api")

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

		io.WriteString(w, "user was successfully created")
	})

	log.Fatalln(http.ListenAndServe(":"+cfg.Http.Port, r))
	
```

Licensed by MIT
