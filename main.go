package main

import (
	"flag"
	"log"
	"net/http"

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
	// todo: check if library.json, if exists -> read from that
}

func handleRequests() {
	mux := http.NewServeMux()
	handler := middleware.CustomMiddleware(mux)
	routes.SetUpRequestHandlers(mux)
	log.Println("desuplayer v2 now listening on http://127.0.0.1:" + serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, handler))
}
