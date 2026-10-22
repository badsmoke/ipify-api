package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func requestIP(query, forwardedFor, remoteAddr string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
	req.RemoteAddr = remoteAddr
	if forwardedFor != "" {
		req.Header.Set("X-Forwarded-For", forwardedFor)
	}

	response := httptest.NewRecorder()
	GetIP(response, req, nil)
	return response
}

func TestGetIPPlainText(t *testing.T) {
	tests := []struct {
		name         string
		forwardedFor string
		remoteAddr   string
		wantIP       string
	}{
		{
			name:       "uses TCP peer when forwarded header is absent",
			remoteAddr: "198.51.100.20:12345",
			wantIP:     "198.51.100.20",
		},
		{
			name:         "uses first forwarded IP and trims whitespace",
			forwardedFor: " 203.0.113.7, 10.0.0.1",
			remoteAddr:   "10.0.0.2:12345",
			wantIP:       "203.0.113.7",
		},
		{
			name:         "falls back when forwarded IP is invalid",
			forwardedFor: "not-an-ip",
			remoteAddr:   "198.51.100.21:12345",
			wantIP:       "198.51.100.21",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := requestIP("", tt.forwardedFor, tt.remoteAddr)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body = %q", response.Code, http.StatusOK, response.Body.String())
			}
			if got := response.Body.String(); got != tt.wantIP {
				t.Errorf("body = %q, want %q", got, tt.wantIP)
			}
			if got := response.Header().Get("Content-Type"); got != "text/plain" {
				t.Errorf("Content-Type = %q, want %q", got, "text/plain")
			}
		})
	}
}

func TestGetIPJSON(t *testing.T) {
	response := requestIP("format=json", "203.0.113.7", "10.0.0.2:12345")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got, want := response.Body.String(), `{"ip":"203.0.113.7"}`; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

func TestGetIPJSONP(t *testing.T) {
	response := requestIP("format=jsonp&callback=myApp.callback", "203.0.113.7", "10.0.0.2:12345")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got, want := response.Body.String(), `myApp.callback({"ip":"203.0.113.7"});`; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
	if got := response.Header().Get("Content-Type"); got != "application/javascript" {
		t.Errorf("Content-Type = %q, want %q", got, "application/javascript")
	}
}

func TestGetIPRejectsInvalidJSONPCallback(t *testing.T) {
	response := requestIP("format=jsonp&callback=x%3Balert(1)%2F%2F", "203.0.113.7", "10.0.0.2:12345")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	if strings.Contains(response.Body.String(), "alert(1)") {
		t.Errorf("response contains untrusted callback: %q", response.Body.String())
	}
}

func TestGetIPRejectsUnknownRemoteAddress(t *testing.T) {
	response := requestIP("", "not-an-ip", "invalid-remote-address")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
