package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadJsonByPath(path string, cfg any) error {
	conFile, err := os.Open(path)
	if err != nil {

		return fmt.Errorf("could not open config file: %w", err)
	}
	defer conFile.Close()

	return json.NewDecoder(conFile).Decode(cfg)
}

func MustJsonYamlByPath(path string, cfg any) {
	if err := LoadJsonByPath(path, cfg); err != nil {
		panic(err)
	}
}
