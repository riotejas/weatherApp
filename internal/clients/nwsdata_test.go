package clients

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_FetchClientData(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "mock response"}`))
	}))
	defer mockServer.Close()

	resp, err := fetchClientData(mockServer.URL)
	if err != nil {
		t.Fatalf("error making request to svr. Err: %v", err)
	}
	expected := `{"message": "mock response"}`
	if string(resp) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", resp, expected)
	}
}
