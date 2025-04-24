package main

import (
	"fmt"
	"kvlt/envs"
	"kvlt/routes"
	"kvlt/store"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func main() {
	// Load Environment Variables
	envs.LoadEnv()

	envs := envs.Gets()

	// Init Store
	store := store.Get()
	store.EnableAutoPersistence(envs.DbPath)

	jobDropExpiredKeys(store, envs.CleanerTime)

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

func jobDropExpiredKeys(store *store.Store, time string) {
	crontab := cron.New()
	_, err := crontab.AddFunc(time, func() {
		if err := store.CleanExpiredKeys(); err != nil {
			log.Printf("Error cleaning expired keys: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job (jobDropExpiredKeys): %v", err)
	}
	crontab.Start()
}
