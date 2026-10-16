package commands

import (
	"errors"
	"fmt"
	"io"
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
	current_flag, current_output, err := parseRunFlags(arg)
	if err != nil {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Flag\tDescription")
		fmt.Fprintln(w, "----\t-----------")
		fmt.Fprintln(w, "-a\tRun all APIs (default)")
		fmt.Fprintln(w, "-n <name>...\tRun only the named APIs")
		fmt.Fprintln(w, "-t\tTerminal output format")
		w.Flush()
		return err
	}

	return runRun(os.Stdout, current_flag, current_output, arg)
}

func parseRunFlags(arg []string) (flag, output, error) {
	current_flag := all
	current_output := terminal

	if len(arg) <= 2 {
		return current_flag, current_output, nil
	}

	matched, err := regexp.MatchString(`^-([an]?t?)$`, arg[2])
	if err != nil {
		return "", "", errors.New("Error while reading flags")
	}
	if !matched {
		return "", "", fmt.Errorf("unsupported flag used in %s", arg[2])
	}

	if strings.Contains(arg[2], "n") {
		current_flag = names
	}
	return current_flag, current_output, nil
}

func runRun(w io.Writer, current_flag flag, current_output output, arg []string) error {

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
				fmt.Fprintln(w, "Didn't found any api named", name)
			}
		}
	}

	return runApis(w, apisToRun, current_output)

}

func runApis(out io.Writer, apisToCheck []configparser.ApiConfig, current_output output) error {
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
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
	for res := range results {
		fmt.Fprintf(w, "%s\t%v\t%s\t%s\n", res.ApiName, res.Status, res.StatusCode, res.TimeTaken)
	}

	w.Flush()
	return nil
}
