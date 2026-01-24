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
		err := commands.RunList()
		if err != nil {
			log.Fatal(err)
		}
	case commands.Add:
		commands.RunAdd()
	case commands.Remove:
		commands.RunRemove(args[0:])
	case commands.Scope:
		commands.RunScope()
	}
}
