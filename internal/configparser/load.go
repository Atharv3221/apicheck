package configparser

import (
	"errors"
	"os"

	"github.com/Atharv3221/apicheck/internal"
	"github.com/Atharv3221/apicheck/internal/initialization"
	"gopkg.in/yaml.v3"
)

func Load() (Config, error) {
	configPath := internal.GetGlobalPath()
	if initialization.LocalConfigExists() {
		configPath = internal.GetLocalPath()
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, errors.New("Error while reading data from config")
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, errors.New("Error while parsing config")
	}
	return cfg, nil
}
