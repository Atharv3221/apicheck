package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseRunFlags(t *testing.T) {
	tests := []struct {
		args     []string
		wantFlag flag
		wantErr  bool
	}{
		{args: []string{"apicheck", "run"}, wantFlag: all},
		{args: []string{"apicheck", "run", "-a"}, wantFlag: all},
		{args: []string{"apicheck", "run", "-t"}, wantFlag: all},
		{args: []string{"apicheck", "run", "-at"}, wantFlag: all},
		{args: []string{"apicheck", "run", "-n", "x"}, wantFlag: names},
		{args: []string{"apicheck", "run", "-nt", "x"}, wantFlag: names},
		{args: []string{"apicheck", "run", "-x"}, wantErr: true},
		{args: []string{"apicheck", "run", "-an"}, wantErr: true},
		{args: []string{"apicheck", "run", "--all"}, wantErr: true},
		{args: []string{"apicheck", "run", "a"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args[2:], " "), func(t *testing.T) {
			f, o, err := parseRunFlags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if f != tt.wantFlag {
				t.Errorf("flag = %q, want %q", f, tt.wantFlag)
			}
			if o != terminal {
				t.Errorf("output = %q, want %q", o, terminal)
			}
		})
	}
}

type row struct{ status, code string }

// parseRows turns the tabwriter table into name -> row, skipping the header.
func parseRows(t *testing.T, out string) map[string]row {
	t.Helper()
	rows := map[string]row{}
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) != 4 || f[0] == "API" {
			continue
		}
		rows[f[0]] = row{status: f[1], code: f[2]}
	}
	return rows
}

func TestRunRun(t *testing.T) {
	const delay = 200 * time.Millisecond

	config := func(base string) string {
		return fmt.Sprintf(`apis:
  - name: ok1
    url: %[1]s/ok1
    method: GET
  - name: ok2
    url: %[1]s/ok2
    method: GET
  - name: broken
    url: %[1]s/broken
    method: GET
  - name: login
    url: %[1]s/login
    method: POST
    header:
      Authorization: Bearer good
    requestBody: '{"user":"a"}'
`, base)
	}

	tests := []struct {
		name     string
		args     []string
		wantRows map[string]row
		wantOut  string // extra substring expected in output
	}{
		{
			name: "no flag runs all",
			args: []string{"apicheck", "run"},
			wantRows: map[string]row{
				"ok1":    {"PASS", "200"},
				"ok2":    {"PASS", "200"},
				"broken": {"FAIL", "500"},
				"login":  {"PASS", "201"},
			},
		},
		{
			name: "-a runs all",
			args: []string{"apicheck", "run", "-a"},
			wantRows: map[string]row{
				"ok1":    {"PASS", "200"},
				"ok2":    {"PASS", "200"},
				"broken": {"FAIL", "500"},
				"login":  {"PASS", "201"},
			},
		},
		{
			name: "-n runs only named",
			args: []string{"apicheck", "run", "-n", "ok1", "broken", "login"},
			wantRows: map[string]row{
				"ok1":    {"PASS", "200"},
				"broken": {"FAIL", "500"},
				"login":  {"PASS", "201"},
			},
		},
		{
			name:     "-n reports unknown names",
			args:     []string{"apicheck", "run", "-n", "ok2", "missing"},
			wantRows: map[string]row{"ok2": {"PASS", "200"}},
			wantOut:  "Didn't found any api named missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inFlight, peak atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n := inFlight.Add(1)
				defer inFlight.Add(-1)
				for {
					p := peak.Load()
					if n <= p || peak.CompareAndSwap(p, n) {
						break
					}
				}
				time.Sleep(delay)

				switch r.URL.Path {
				case "/ok1", "/ok2":
					w.WriteHeader(http.StatusOK)
				case "/login":
					body, _ := io.ReadAll(r.Body)
					if r.Header.Get("Authorization") == "Bearer good" &&
						r.Header.Get("Content-Type") == "application/json" &&
						string(body) == `{"user":"a"}` {
						w.WriteHeader(http.StatusCreated)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))
			defer srv.Close()

			work := t.TempDir()
			t.Setenv("HOME", t.TempDir())
			t.Chdir(work)
			if err := os.WriteFile(filepath.Join(work, ".apicheck.yaml"), []byte(config(srv.URL)), 0644); err != nil {
				t.Fatal(err)
			}

			f, o, err := parseRunFlags(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			start := time.Now()
			if err := runRun(&buf, f, o, tt.args); err != nil {
				t.Fatalf("runRun() error = %v", err)
			}
			elapsed := time.Since(start)
			out := buf.String()

			got := parseRows(t, out)
			if len(got) != len(tt.wantRows) {
				t.Errorf("got %d rows, want %d:\n%s", len(got), len(tt.wantRows), out)
			}
			for name, want := range tt.wantRows {
				if got[name] != want {
					t.Errorf("row %s = %+v, want %+v", name, got[name], want)
				}
			}
			if tt.wantOut != "" && !strings.Contains(out, tt.wantOut) {
				t.Errorf("output missing %q:\n%s", tt.wantOut, out)
			}

			// Parallelism: every selected API must be in flight at once, and
			// the whole run must take about one delay, not one per API.
			if int(peak.Load()) != len(tt.wantRows) {
				t.Errorf("peak concurrent requests = %d, want %d", peak.Load(), len(tt.wantRows))
			}
			if max := 2 * delay; len(tt.wantRows) > 1 && elapsed >= max {
				t.Errorf("run took %v, want < %v (requests not parallel)", elapsed, max)
			}
		})
	}
}
