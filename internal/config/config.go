package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultHTTPPort               = "8080"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
)

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Host               string        `yaml:"host"`
	Port               string        `yaml:"port"`
	ReadTimeout        time.Duration `yaml:"readTimeout"`
	WriteTimeout       time.Duration `yaml:"writeTimeout"`
	MaxHeaderMegabytes int           `yaml:"maxHeaderBytes"`
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Port:               defaultHTTPPort,
			ReadTimeout:        defaultHTTPRWTimeout,
			WriteTimeout:       defaultHTTPRWTimeout,
			MaxHeaderMegabytes: defaultHTTPMaxHeaderMegabytes,
		},
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
