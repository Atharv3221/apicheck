package main

import (
	"fmt"
	"os"

	"github.com/Atharv3221/apicheck/internal/commands"
	"github.com/Atharv3221/apicheck/internal/initialization"
)

const usage = `Usage: apicheck <command> [arguments]

Commands:
  init                    Create a local .apicheck.yaml in the current directory
  add                     Add an API interactively
  list                    List configured APIs
  run [-a | -n <name>...] Run all checks (-a, default) or only the named ones (-n)
  remove <name>...        Remove APIs by name
  scope                   Show whether the local or global config is in use`

func init() {
	// set up
	initialization.Initialize()
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	switch commands.Command(args[1]) {
	case commands.Init:
		if err := initialization.CreateLocalConfig(); err != nil {
			fail(err)
		}
	case commands.Run:
		if err := commands.RunRun(args); err != nil {
			fail(err)
		}
	case commands.List:
		if err := commands.RunList(); err != nil {
			fail(err)
		}
	case commands.Add:
		if err := commands.RunAdd(); err != nil {
			fail(err)
		}
	case commands.Remove:
		if err := commands.RunRemove(args); err != nil {
			fail(err)
		}
	case commands.Scope:
		commands.RunScope()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n%s\n", args[1], usage)
		os.Exit(1)
	}
}
