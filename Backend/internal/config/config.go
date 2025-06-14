package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Configuration struct {
	ParameterName        string
	ClusterEndpoint      string
	Port                 string
	User                 string
	DatabaseName         string
	Region               string
	TokenRefreshInterval int
}

func LoadConfiguration() (*Configuration, error) {
	var TokenRI = 900
	tri, err := strconv.Atoi(os.Getenv("TOKEN_REFRESH_INTERVAL"))
	if err == nil {
		TokenRI = tri
	}
	cfg := &Configuration{
		ParameterName:        os.Getenv("SSM_PARAMETER_NAME"),
		ClusterEndpoint:      os.Getenv("CLUSTER_ENDPOINT"),
		Port:                 os.Getenv("DB_PORT"),
		User:                 os.Getenv("CLUSTER_USER"),
		DatabaseName:         os.Getenv("DB_NAME"),
		Region:               os.Getenv("REGION"),
		TokenRefreshInterval: TokenRI,
	}
	return cfg, nil
}
