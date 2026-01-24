package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Atharv3221/apicheck/internal/commands"
	"github.com/Atharv3221/apicheck/internal/initialization"
)

func init() {
	// set up
	initialization.Initialize()
}

func main() {
	args := os.Args

	switch commands.Command(args[1]) {
	case commands.Init:
		err := initialization.CreateLocalConfig()
		if err != nil {
			log.Fatal("Problem while creating local config")
		}
	case commands.Run:
		fmt.Println("Command Run")
	case commands.List:
		fmt.Println("command List")
	case commands.Remove:
		fmt.Println("command Remove")
	case commands.Scope:
		fmt.Println("command Scope")
	}
}
