package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ComputingPower   int    `envconfig:"COMPUTING_POWER"`
	Host             string `envconfig:"HOST" default:"localhost"`
	OrchestratorPort string `envconfig:"ORCHESTRATOR_GRPC_PORT"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg Config
	err := envconfig.Process("agent_service", &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
