package main

import (
	"fmt"
	"kvlt/envs"
	"kvlt/routes"
	"kvlt/store"
	"log"
	"net/http"
)

func main() {
	// Load Environment Variables
	envs.LoadEnv()

	envs := envs.Gets()

	// Init Store
	store := store.Get()
	store.EnableAutoPersistence(envs.DbPath)

	startServer(envs.Host, envs.Port)
}

func startServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	router := routes.SetupRouter()

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Server start on: http://%s", addr)
	log.Fatal(server.ListenAndServe())
}
