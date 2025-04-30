package envs

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Envs struct {
	IntervalAnalyzeBuffer time.Duration `env:"INTERVAL_ANALYZE_BUFFER,required"`
	IntervalSnapshot      time.Duration `env:"INTERVAL_SNAPSHOT,required"`
	IntervalCleaner       time.Duration `env:"INTERVAL_CLEANER,required"`
	QuantityBuffer        int           `env:"QUANTITY_BUFFER,required"`
	AofFolderPath         string        `env:"AOF_FOLDER_PATH,required"`
	SnapshotFolderPath    string        `env:"SNAPSHOT_FOLDER_PATH,required"`
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}
}

func Gets() Envs {
	var envs Envs

	if err := env.Parse(&envs); err != nil {
		fmt.Printf("Error parsing env variables: %v\n", err)
		os.Exit(1)
	}

	return envs
}
