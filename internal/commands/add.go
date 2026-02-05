package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

func runAdd(api configparser.ApiConfig) error {
	cfg, err := configparser.Load()
	if err != nil {
		return err
	}
	cfg.Apis = append(cfg.Apis, api)
	configparser.Write(&cfg)
	return nil
}

func RunAdd() error {
	reader := bufio.NewReader(os.Stdin)
	api := configparser.ApiConfig{
		Header: make(map[string]string),
	}

	fmt.Println("─── Adding New API ───")

	readInput := func(prompt string) string {
		fmt.Printf("%-15s: ", prompt)
		val, _ := reader.ReadString('\n')
		return strings.TrimSpace(val)
	}

	api.Name = readInput("API Name")
	api.Url = readInput("URL")
	api.Method = strings.ToUpper(readInput("Method"))

	fmt.Println("Enter Header (Format: Key: Value, or leave empty):")
	headerRaw := readInput("Header")
	if headerRaw != "" {
		parts := strings.SplitN(headerRaw, ":", 2)
		if len(parts) == 2 {
			api.Header[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	fmt.Println("Request Body (Press Enter if empty):")
	api.RequestBody = readInput(">")

	err := runAdd(api)
	if err != nil {
		return fmt.Errorf("failed to save API: %w", err)
	}

	fmt.Println("✔ API added successfully!")
	return nil
}
