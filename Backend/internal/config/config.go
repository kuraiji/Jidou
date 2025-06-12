package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Configuration struct {
	ParameterName string
}

func LoadConfiguration() (*Configuration, error) {
	cfg := &Configuration{
		ParameterName: os.Getenv("SSM_PARAMETER_NAME"),
	}
	return cfg, nil
}
