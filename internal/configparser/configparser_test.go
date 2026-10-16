package configparser

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const configFile = ".apicheck.yaml"

// setupDirs points HOME and the working directory at fresh temp dirs so the
// global (~/.apicheck.yaml) and local (./.apicheck.yaml) configs are isolated.
func setupDirs(t *testing.T) (home, work string) {
	t.Helper()
	home, work = t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(work)
	return home, work
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		global  string // empty = file not created
		local   string // empty = file not created
		want    Config
		wantErr bool
	}{
		{
			name: "global only",
			global: `apis:
  - name: g
    url: http://global
    method: GET
`,
			want: Config{Apis: []ApiConfig{{Name: "g", Url: "http://global", Method: "GET"}}},
		},
		{
			name: "local overrides global",
			global: `apis:
  - name: g
    url: http://global
    method: GET
`,
			local: `apis:
  - name: l
    url: http://local
    method: POST
`,
			want: Config{Apis: []ApiConfig{{Name: "l", Url: "http://local", Method: "POST"}}},
		},
		{
			name: "headers and request body",
			global: `apis:
  - name: login
    url: http://example.com/login
    method: POST
    header:
      Authorization: Bearer abc
      X-Trace: "1"
    requestBody: '{"user":"a"}'
`,
			want: Config{Apis: []ApiConfig{{
				Name:        "login",
				Url:         "http://example.com/login",
				Method:      "POST",
				Header:      map[string]string{"Authorization": "Bearer abc", "X-Trace": "1"},
				RequestBody: `{"user":"a"}`,
			}}},
		},
		{
			name:   "default template has no apis",
			global: "# Default configuration for apicheck\n\napis:\n",
			want:   Config{},
		},
		{
			name:    "malformed yaml",
			global:  "apis: [unclosed",
			wantErr: true,
		},
		{
			name:    "missing config",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, work := setupDirs(t)
			if tt.global != "" {
				writeFile(t, filepath.Join(home, configFile), tt.global)
			}
			if tt.local != "" {
				writeFile(t, filepath.Join(work, configFile), tt.local)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	cfg := Config{Apis: []ApiConfig{
		{Name: "plain", Url: "http://a", Method: "GET"},
		{
			Name:        "full",
			Url:         "http://b",
			Method:      "POST",
			Header:      map[string]string{"Authorization": "Bearer x"},
			RequestBody: `{"k":"v"}`,
		},
	}}

	tests := []struct {
		name        string
		localExists bool
		readOnly    bool
		wantErr     bool
	}{
		{name: "writes global when no local config"},
		{name: "writes local when local config exists", localExists: true},
		{name: "read-only config returns error", readOnly: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, work := setupDirs(t)
			globalPath := filepath.Join(home, configFile)
			localPath := filepath.Join(work, configFile)
			writeFile(t, globalPath, "apis:\n")
			if tt.localExists {
				writeFile(t, localPath, "apis:\n")
			}
			if tt.readOnly {
				if os.Geteuid() == 0 {
					t.Skip("root ignores file permissions")
				}
				if err := os.Chmod(globalPath, 0444); err != nil {
					t.Fatal(err)
				}
			}

			err := Write(&cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Round trip: what was written must load back identically.
			got, err := Load()
			if err != nil {
				t.Fatalf("Load() after Write() error = %v", err)
			}
			if !reflect.DeepEqual(got, cfg) {
				t.Errorf("round trip = %+v, want %+v", got, cfg)
			}

			// The other scope's file must be untouched.
			untouched := localPath
			if tt.localExists {
				untouched = globalPath
			}
			if data, err := os.ReadFile(untouched); err == nil && string(data) != "apis:\n" {
				t.Errorf("%s was modified: %q", untouched, data)
			}
		})
	}
}
