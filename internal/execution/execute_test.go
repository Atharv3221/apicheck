package execution

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Atharv3221/apicheck/internal/configparser"
)

type received struct {
	method string
	header http.Header
	body   string
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		header     map[string]string
		body       string
		respStatus int

		wantStatus      Result
		wantStatusCode  string
		wantHeaders     map[string]string // headers the server must see
		wantBody        string
		wantContentType string // "" = must not be set
	}{
		{
			name:           "GET 200 passes",
			method:         "GET",
			respStatus:     http.StatusOK,
			wantStatus:     PASS,
			wantStatusCode: "200",
		},
		{
			name:            "POST sends JSON body with default content type",
			method:          "POST",
			body:            `{"username":"emilys","password":"emilyspass"}`,
			respStatus:      http.StatusCreated,
			wantStatus:      PASS,
			wantStatusCode:  "201",
			wantBody:        `{"username":"emilys","password":"emilyspass"}`,
			wantContentType: "application/json",
		},
		{
			name:   "custom headers reach server",
			method: "GET",
			header: map[string]string{
				"Authorization": "Bearer secret-token",
				"X-Api-Key":     "k123",
			},
			respStatus:     http.StatusOK,
			wantStatus:     PASS,
			wantStatusCode: "200",
			wantHeaders: map[string]string{
				"Authorization": "Bearer secret-token",
				"X-Api-Key":     "k123",
			},
		},
		{
			name:            "explicit content type is not overridden",
			method:          "PUT",
			header:          map[string]string{"Content-Type": "text/plain"},
			body:            "hello",
			respStatus:      http.StatusOK,
			wantStatus:      PASS,
			wantStatusCode:  "200",
			wantBody:        "hello",
			wantContentType: "text/plain",
		},
		{
			name:            "headers and body together",
			method:          "POST",
			header:          map[string]string{"Authorization": "Bearer t"},
			body:            `{"name":"x"}`,
			respStatus:      http.StatusCreated,
			wantStatus:      PASS,
			wantStatusCode:  "201",
			wantHeaders:     map[string]string{"Authorization": "Bearer t"},
			wantBody:        `{"name":"x"}`,
			wantContentType: "application/json",
		},
		{
			name:           "lowercase method is normalised",
			method:         "delete",
			respStatus:     http.StatusNoContent,
			wantStatus:     PASS,
			wantStatusCode: "204",
		},
		{
			name:            "401 fails",
			method:          "POST",
			body:            `{}`,
			respStatus:      http.StatusUnauthorized,
			wantStatus:      FAIL,
			wantStatusCode:  "401",
			wantBody:        `{}`,
			wantContentType: "application/json",
		},
		{
			name:           "404 fails",
			method:         "GET",
			respStatus:     http.StatusNotFound,
			wantStatus:     FAIL,
			wantStatusCode: "404",
		},
		{
			name:           "500 fails",
			method:         "GET",
			respStatus:     http.StatusInternalServerError,
			wantStatus:     FAIL,
			wantStatusCode: "500",
		},
		{
			name:           "3xx is not a pass",
			method:         "GET",
			respStatus:     http.StatusNotModified,
			wantStatus:     FAIL,
			wantStatusCode: "304",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make(chan received, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				got <- received{method: r.Method, header: r.Header.Clone(), body: string(b)}
				w.WriteHeader(tt.respStatus)
			}))
			defer srv.Close()

			out := Execute(configparser.ApiConfig{
				Name:        tt.name,
				Url:         srv.URL,
				Method:      tt.method,
				Header:      tt.header,
				RequestBody: tt.body,
			}, "t")

			if out.ApiName != tt.name {
				t.Errorf("ApiName = %q, want %q", out.ApiName, tt.name)
			}
			if out.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", out.Status, tt.wantStatus)
			}
			if out.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode = %q, want %q", out.StatusCode, tt.wantStatusCode)
			}
			if out.TimeTaken == "" {
				t.Error("TimeTaken is empty")
			}

			req := <-got
			if want := strings.ToUpper(tt.method); req.method != want {
				t.Errorf("server saw method %q, want %q", req.method, want)
			}
			for k, v := range tt.wantHeaders {
				if g := req.header.Get(k); g != v {
					t.Errorf("server saw header %s = %q, want %q", k, g, v)
				}
			}
			if req.body != tt.wantBody {
				t.Errorf("server saw body %q, want %q", req.body, tt.wantBody)
			}
			if g := req.header.Get("Content-Type"); g != tt.wantContentType {
				t.Errorf("server saw Content-Type %q, want %q", g, tt.wantContentType)
			}
		})
	}
}

func TestExecuteRequestErrors(t *testing.T) {
	closed := httptest.NewServer(http.NotFoundHandler())
	closedURL := closed.URL
	closed.Close()

	tests := []struct {
		name string
		url  string
	}{
		{name: "invalid url", url: "://not-a-url"},
		{name: "connection refused", url: closedURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := Execute(configparser.ApiConfig{Name: tt.name, Url: tt.url, Method: "GET"}, "t")
			if out.Status != FAIL {
				t.Errorf("Status = %v, want FAIL", out.Status)
			}
			if out.StatusCode != "-" {
				t.Errorf("StatusCode = %q, want \"-\"", out.StatusCode)
			}
		})
	}
}
