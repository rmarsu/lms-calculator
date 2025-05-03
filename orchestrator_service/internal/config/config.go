package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrpcPort             string `envconfig:"ORCHESTRATOR_GRPC_PORT"`
	RestPort             string `envconfig:"ORCHESTRATOR_REST_PORT"`
	SqlitePath           string `envconfig:"ORCHESTRATOR_SQLITE_PATH"`
	JwtSecret            string `envconfig:"JWT_SECRET"`
	TimeAdditionMs       int64  `envconfig:"TIME_ADDITION_MS"`
	TimeSubtractionMs    int64  `envconfig:"TIME_SUBTRACTION_MS"`
	TimeMultiplicationMs int64  `envconfig:"TIME_MULTIPLICATION_MS"`
	TimeDivisionMs       int64  `envconfig:"TIME_DIVISION_MS"`
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
