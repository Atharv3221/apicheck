package execution

type Result string

const (
	PASS Result = "PASS"
	FAIL Result = "FAIL"
)

type TerminalOutput struct {
	Status     Result
	StatusCode string
	TimeTaken  string
	ApiName    string
}
