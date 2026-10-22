package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	newHandler().ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
}

func TestLogIPsEnabledByEnvironment(t *testing.T) {
	t.Setenv("LOG_IPS", "true")
	logs := captureLogs(t)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "198.51.100.20:12345"
	response := httptest.NewRecorder()

	newHandler().ServeHTTP(response, request)

	if !strings.Contains(logs.String(), "client_ip=198.51.100.20") {
		t.Errorf("logs = %q, want client IP", logs.String())
	}
}

func TestLogIPsDisabledByDefault(t *testing.T) {
	t.Setenv("LOG_IPS", "")
	logs := captureLogs(t)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "198.51.100.20:12345"
	response := httptest.NewRecorder()

	newHandler().ServeHTTP(response, request)

	if got := logs.String(); got != "" {
		t.Errorf("logs = %q, want no IP log", got)
	}
}

func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()

	logs := &bytes.Buffer{}
	previous := log.Writer()
	log.SetOutput(logs)
	t.Cleanup(func() { log.SetOutput(previous) })
	return logs
}
