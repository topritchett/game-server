package main

import (
	"log"
	"net/http"

	"github.com/topritchett/game-server/server"
)

func main() {
	mux := http.NewServeMux()
	server.New(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
