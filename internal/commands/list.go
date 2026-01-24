package commands

import (
	"fmt"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

func RunList() error {
	cfg, err := configparser.Load()
	if err != nil {
		return err
	}
	fmt.Println("API list loaded!")
	for _, api := range cfg.Apis {
		fmt.Println("Name: " + api.Name)
		fmt.Println("URL: " + api.Url)
		fmt.Println("Method: " + api.Method)
		fmt.Println("")
	}
	return nil
}
