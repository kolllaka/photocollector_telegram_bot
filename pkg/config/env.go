package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	DEFAULT_ENV_PATH = ".env"
)

func LoadEnvByPath(path string, cfg any) error {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return err
	}

	return nil
}

func MustLoadEnvByPath(path string, cfg any) {
	if err := LoadEnvByPath(path, cfg); err != nil {
		panic(err)
	}
}

func LoadEnv(cfg any) error {
	return LoadEnvByPath(DEFAULT_ENV_PATH, cfg)
}

func MustLoadEnv(cfg any) {
	if err := LoadEnv(cfg); err != nil {
		panic(err)
	}
}
