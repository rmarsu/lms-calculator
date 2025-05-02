package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrpcPort   string `envconfig:"ORCHESTRATOR_GRPC_PORT"`
	RestPort   string `envconfig:"ORCHESTRATOR_REST_PORT"`
	SqlitePath string `envconfig:"ORCHESTRATOR_SQLITE_PATH"`
	JwtSecret  string `envconfig:"JWT_SECRET"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg Config
	err := envconfig.Process("orchestrator_service", &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
