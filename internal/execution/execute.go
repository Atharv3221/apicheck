package execution

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Atharv3221/apicheck/internal/configparser"
	"github.com/Atharv3221/apicheck/internal/constants/methods"
)

func Execute(api configparser.ApiConfig, output string) TerminalOutput {
	terminal_output := TerminalOutput{}
	terminal_output.ApiName = api.Name
	api_method := methods.Method(strings.ToUpper(api.Method))
	start := time.Now()
	res, err := http.NewRequest(string(api_method), api.Url, nil)
	time_duration := time.Since(start)
	durationStr := time_duration.Round(time.Millisecond).String()
	if err != nil {
		terminal_output.Status = FAIL
		terminal_output.TimeTaken = durationStr
		terminal_output.StatusCode = ""
		return terminal_output
	} else {
		defer res.Body.Close()
		terminal_output.StatusCode = strconv.Itoa(res.Response.StatusCode)
		terminal_output.TimeTaken = durationStr

		if res.Response.StatusCode >= 200 && res.Response.StatusCode < 300 {
			terminal_output.Status = PASS
		} else {
			terminal_output.Status = FAIL
		}

	}
	return terminal_output
}
