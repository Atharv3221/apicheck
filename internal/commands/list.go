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
	fmt.Println("")

	for _, api := range cfg.Apis {
		fmt.Println("Name:   " + api.Name)
		fmt.Println("URL:    " + api.Url)
		fmt.Println("Method: " + api.Method)

		if len(api.Header) > 0 {
			for key, value := range api.Header {
				fmt.Printf("Header: %s: %s\n", key, value)
			}
		} else {
			fmt.Println("Header: None")
		}

		if api.RequestBody != "" {
			fmt.Println("Body:   " + api.RequestBody)
		}

		fmt.Println("")
	}
	return nil
}
