package config

import (
	"log"
	"os"
)

// Config provides DSN and Port
type Config struct {
	DSN  string
	Port string
}

// Load provides port for server and link to DB from .env
func Load() *Config {
	config := Config{}

	config.Port = os.Getenv("SUBSCRIPTION_PORT")
	if config.Port == "" {
		log.Fatal("SUBSCRIPTION_PORT is not set in env")
	}
	config.DSN = os.Getenv("DATABASE_URL")
	if config.DSN == "" {
		log.Fatal("DATABASE_URL is not set in env")
	}
	return &config

}
