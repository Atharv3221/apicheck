package initialization

import (
	"errors"
	"fmt"
	"os"

	"github.com/Atharv3221/apicheck/internal"
)

func CreateLocalConfig() error {
	file := internal.GetLocalPath()
	fmt.Println("Building .apicheck.yaml for current directory")
	if LocalConfigExists() {
		fmt.Println("Local config already exists")
	} else {
		err := os.WriteFile(file, template, 0644)
		if err != nil {
			return errors.New("Problem while creating local config")
		}
	}
	return nil
}

func LocalConfigExists() bool {
	file := internal.GetLocalPath()
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}
