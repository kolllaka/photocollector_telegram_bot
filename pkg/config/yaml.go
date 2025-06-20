package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYamlByPath(path string, cfg any) error {
	conFile, err := os.Open(path)
	if err != nil {

		return fmt.Errorf("could not open config file: %w", err)
	}
	defer conFile.Close()

	return yaml.NewDecoder(conFile).Decode(cfg)
}

func MustLoadYamlByPath(path string, cfg any) {
	if err := LoadYamlByPath(path, cfg); err != nil {
		panic(err)
	}
}
