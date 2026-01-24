package initialization

import (
	_ "embed"
	"errors"
	"os"

	"github.com/Atharv3221/apicheck/internal"
)

//go:embed default_config.yaml
var template []byte

// check if global config file exists
func globalConfigExists() bool {
	file := internal.GetGlobalPath()
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}

// create global config
func createGlobalConfig() error {
	file := internal.GetGlobalPath()
	err := os.WriteFile(file, template, 0644)
	if err != nil {
		return errors.New("Problem while creating global config")
	}
	return nil
}
