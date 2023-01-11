package main

import (
	"log"
	"net/http"

	"github.com/topritchett/game-server/api"
	"github.com/topritchett/game-server/server"
)

func main() {
	webMux := http.NewServeMux()
	server.New(webMux)
	go func() {
		log.Println("Web server started on :8000")
		log.Fatal(http.ListenAndServe(":8000", webMux))
	}()

	apiMux := http.NewServeMux()
	api.New(apiMux)
	go func() {
		log.Println("API server started on :9000")
		log.Fatal(http.ListenAndServe(":9000", apiMux))
	}()

	log.Println("Server started")
	select {}
}
