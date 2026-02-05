package commands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/Atharv3221/apicheck/internal/configparser"
	"github.com/Atharv3221/apicheck/internal/execution"
)

type flag string

const (
	all   flag = "a"
	names flag = "n"
)

type output string

const (
	terminal output = "t"
)

func RunRun(arg []string) error {
	current_flag := all
	current_output := terminal

	matched := false
	var err error
	if (len(arg)) > 2 {
		pattern := `^-([an]?t?)$`
		matched, err = regexp.MatchString(pattern, arg[2])
		if err != nil {
			return errors.New("Error while reading flags")
		}
	}

	if matched {
		if strings.Contains(arg[2], "a") {
			current_flag = all
		} else if strings.Contains(arg[2], "n") {
			current_flag = names
		}
		if strings.Contains(arg[2], "t") {
			current_output = terminal
		}

	} else if len(arg) > 2 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Flag\tDescription")
		fmt.Fprintln(w, "----\t-----------")
		fmt.Fprintln(w, "-a\tShow all details")
		fmt.Fprintln(w, "-n\tShow names only")
		fmt.Fprintln(w, "-t\tTerminal output format")
		w.Flush()
		return fmt.Errorf("unsupported flag used in %s", arg[2])
	}

	return runRun(current_flag, current_output, arg)
}

func runRun(current_flag flag, current_output output, arg []string) error {

	cfg, err := configparser.Load()
	if err != nil {
		return err
	}

	var apisToRun []configparser.ApiConfig

	if current_flag == all {
		apisToRun = cfg.Apis
	} else if current_flag == names && len(arg) > 3 {
		for _, name := range arg[3:] {
			found := false
			for _, api := range cfg.Apis {
				if name == api.Name {
					apisToRun = append(apisToRun, api)
					found = true
					break
				}
			}
			if !found {
				fmt.Println("Didn't found any api named", name)
			}
		}
	}

	return runApis(apisToRun, current_output)

}

func runApis(apisToCheck []configparser.ApiConfig, current_output output) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "API NAME\tSTATUS\tSTATUS_CODE\tTIME")

	results := make(chan execution.TerminalOutput, len(apisToCheck))
	var wg sync.WaitGroup

	for _, api := range apisToCheck {
		wg.Add(1)
		go func(a configparser.ApiConfig) {
			defer wg.Done()
			results <- execution.Execute(a, string(current_output))
		}(api)
	}

	// Close channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	for out := range results {
		fmt.Fprintf(w, "%s\t%v\t%s\t%s\n", out.ApiName, out.Status, out.StatusCode, out.TimeTaken)
	}

	w.Flush()
	return nil
}
