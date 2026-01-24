package configparser

import (
	"errors"
	"os"

	"github.com/Atharv3221/apicheck/internal"
	"github.com/Atharv3221/apicheck/internal/initialization"
	"gopkg.in/yaml.v3"
)

func Write(cfg *Config) error {
	config_data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	configPath := internal.GetGlobalPath()
	if initialization.LocalConfigExists() {
		configPath = internal.GetLocalPath()
	}
	err = os.WriteFile(configPath, config_data, 0644)
	if err != nil {
		return errors.New("Error while adding api")
	}
	return nil
}
