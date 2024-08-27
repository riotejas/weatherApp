package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HealthHandler(t *testing.T) {
	s := &Server{}
	svr := httptest.NewServer(http.HandlerFunc(s.HealthHandler))
	defer svr.Close()
	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Fatalf("error making request to svr. Err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"App is healthy\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
