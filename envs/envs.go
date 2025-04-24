package envs

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Envs struct {
	Port             int    `env:"PORT,required"`
	Host             string `env:"HOST,required"`
	SnapshotTime     string `env:"SNAPSHOT_TIME,required"`
	SnapshotFileName string `env:"SNAPSHOT_FILENAME"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Erreur lors du chargement du fichier .env: %v\n", err)
		os.Exit(1)
	}
}

func Gets() Envs {
	var envs Envs

	if err := env.Parse(&envs); err != nil {
		fmt.Printf("Erreur lors du parsing des variables d'environnement: %v\n", err)
		os.Exit(1)
	}

	return envs
}
