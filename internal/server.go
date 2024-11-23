package internal

import (
	config "gomail/internal/config"
	routes "gomail/internal/routes"
	"log"
	"net/http"
)

func Run() {
	// Load the env file
	config.LoadEnv()

	mux := http.NewServeMux()

	routes.RegisterRoutes(mux)

	cors := config.Cors()

	server := &http.Server{
		Addr:    ":8000",
		Handler: cors.Handler(mux),
	}

	log.Printf("Server started: http://localhost%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
