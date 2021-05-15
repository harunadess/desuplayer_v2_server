package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/middleware"
	"github.com/jordanjohnston/desuplayer_v2/routes"
)

var serverPort string

func init() {
	port := flag.String("port", "4444", "specify a port to run the server on")
	flag.Parse()
	serverPort = *port
}

func main() {
	handleRequests()
}

func handleRequests() {
	mux := http.NewServeMux()
	handler := middleware.CustomMiddleware(mux)
	routes.SetUpRequestHandlers(mux)
	library.LoadLibrary()
	log.Println("desuplayer v2 now listening on http://localhost:" + serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, handler))
}
