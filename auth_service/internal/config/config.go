package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrpcPort   string `envconfig:"GRPC_PORT"`
	RestPort   string `envconfig:"REST_PORT"`
	SqlitePath string `envconfig:"SQLITE_PATH"`
	JwtSecret  string `envconfig:"JWT_SECRET"`
	HasherSalt string `envconfig:"HASHER_SALT"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var cfg Config
	err := envconfig.Process("auth_service", &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
