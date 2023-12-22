# Root
Simple HTTP router for Golang
```go
func main() {

	r := api.NewRouter("/api")

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	server := api.NewServer(cfg.Http.Port)

	server.Router("/api", &api.Handler{Router: r})

	r.GET("/users", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("GETRT USERS")
	})
	r.POST("/users", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("POST USERS")
	})
	r.POST("/users/:id", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("POST USERS :ID")
	})
	r.GET("/users/:id/contacts", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("GET USERS :ID CONTACTS")
		io.WriteString(w, "contact posted [OK]")
	})
	r.POST("/users/:id/contacts", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {})

	log.Printf("Server started on port %s", cfg.Http.Port)

	log.Fatal(server.Run())
}
```
