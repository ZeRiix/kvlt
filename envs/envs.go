package envs

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Envs struct {
	Port   int    `env:"PORT,required"`
	Host   string `env:"HOST,required"`
	DbPath string `env:"DB_PATH" envDefault:"./db/store.json"`
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
