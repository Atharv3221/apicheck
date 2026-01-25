package execution

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

type method string

const (
	get  method = "GET"
	post method = "POST"
)

func Execute(api configparser.ApiConfig, output string) TerminalOutput {
	terminal_output := TerminalOutput{}
	terminal_output.ApiName = api.Name
	switch method(strings.ToUpper(api.Method)) {
	case get:
		start := time.Now()
		res, err := http.Get(api.Url)
		time_duration := time.Since(start)
		durationStr := time_duration.Round(time.Millisecond).String()

		if err != nil {
			terminal_output.Status = FAIL
			terminal_output.TimeTaken = durationStr
			terminal_output.StatusCode = ""
			return terminal_output
		} else {
			defer res.Body.Close()
			terminal_output.StatusCode = strconv.Itoa(res.StatusCode)
			terminal_output.TimeTaken = durationStr // Use rounded string here too

			if res.StatusCode >= 200 && res.StatusCode < 300 {
				terminal_output.Status = PASS
			} else {
				terminal_output.Status = FAIL
			}
		}
	default:
		fmt.Println("Running method", api.Method)

	}
	return terminal_output
}
