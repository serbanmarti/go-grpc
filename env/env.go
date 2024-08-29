package env

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Conf struct {
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	Port        int    `env:"PORT" envDefault:"8080"`
	TokenSecret string `env:"SECRET_TOKEN,required"`
	TokenHeader string `env:"TOKEN_HEADER" envDefault:"x-auth-token"`
}

var lock = &sync.Mutex{}
var envInstance *Conf

// GetEnvironment returns the environment configuration (singleton)
func GetEnvironment() *Conf {
	lock.Lock()
	defer lock.Unlock()
	if envInstance == nil {
		// Load the .env file, if it exists
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		// Parse the environment variables
		envInstance = &Conf{}
		if err := env.Parse(envInstance); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}

	return envInstance
}
