package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

func runAdd(api configparser.ApiConfig) error {
	cfg, err := configparser.Load()
	if err != nil {
		return err
	}
	cfg.Apis = append(cfg.Apis, api)
	configparser.Write(cfg)
	return nil
}

func RunAdd() error {
	reader := bufio.NewReader(os.Stdin)
	api := new(configparser.ApiConfig)
	fmt.Println("Adding new API")
	fmt.Println("Enter name for API: ")
	fmt.Fscan(reader, &api.Name)
	fmt.Println("Enter api URL")
	fmt.Fscan(reader, &api.Url)
	fmt.Println("Enter method: ")
	fmt.Fscan(reader, &api.Method)
	err := runAdd(*api)
	if err != nil {
		return err
	}
	return nil
}
