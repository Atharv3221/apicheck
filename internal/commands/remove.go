package commands

import (
	"fmt"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

func RunRemove(args []string) error {
	cfg, err := configparser.Load()
	if err != nil {
		return err
	}

	if len(args) < 3 {
		fmt.Println("Please enter names of api to be removed!")
		return nil
	}

	toDelete := make([]bool, len(cfg.Apis))

	for i := 2; i < len(args); i++ {

		found := false
		for j, api := range cfg.Apis {

			if api.Name == args[i] {
				toDelete[j] = true
				found = true
			}
		}
		if !found {
			fmt.Printf("Can't find api with name: %s\n", args[i])
		}
	}

	var updatedApis = new(configparser.Config)
	for i, api := range cfg.Apis {
		if !toDelete[i] {
			updatedApis.Apis = append(updatedApis.Apis, api)
		}
	}
	configparser.Write(updatedApis)
	return nil
}
