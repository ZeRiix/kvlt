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
	envs.LoadEnv()

	envs := envs.Gets()

	initializeStore(envs.SnapshotFileName)

	setupCronSnapshot(envs.SnapshotTime)

	startServer(envs.Host, envs.Port)
}

func initializeStore(snapshotName string) {
	store.Get()

	var snapshotToLoad string
	var err error

	if snapshotName != "" {
		snapshotToLoad = snapshotName
		log.Printf("Utilisation du snapshot spécifié: %s", snapshotName)
	} else {
		snapshotToLoad, err = store.GetLatestSnapshot()
		if err != nil {
			log.Fatalf("Erreur lors de la récupération du dernier snapshot: %v", err)
		}
	}

	if snapshotToLoad != "" {
		if err := store.LoadSnapshot(snapshotToLoad); err != nil {
			log.Fatalf("Erreur lors du chargement du snapshot %s: %v", snapshotToLoad, err)
		}
		log.Printf("Snapshot chargé: %s", snapshotToLoad)
	} else {
		log.Println("Aucun snapshot trouvé ou spécifié.")
	}

	log.Println("Store initialisé.")
}

func setupCronSnapshot(time string) {
	crontab := cron.New()
	_, err := crontab.AddFunc(time, func() {
		if err := store.MakeSnapshot(); err != nil {
			log.Printf("Erreur lors de la création du snapshot: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Erreur lors de la configuration du cron job: %v", err)
	}
	crontab.Start()
}

func startServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	router := routes.SetupRouter()

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Serveur démarré sur http://%s", addr)
	log.Fatal(server.ListenAndServe())
}
