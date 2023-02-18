package config

import (
	"os"

	"github.com/joho/godotenv"
)

type HealthConfig struct {
	ClientID     string
	ClientSecret string
}

func LoadConfig() (*HealthConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	conf := HealthConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
	return &conf, err
}
