package config

import (
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress string        `yaml:"server_address"`
	AddressDB     string        `yaml:"db_string"`
	ServerTimeout time.Duration `yaml:"server_timeout"`
}

func CreateConfig(log *slog.Logger, path string) (*Config, error) {
	
	f, err := os.Open(path)
	if err != nil {
		log.Error("CreateConfig","os.Open", err)
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}

	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		log.Error("CreateConfig", "yaml.Decode", err)
		return nil, err
	}

	return cfg, nil
}
