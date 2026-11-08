package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealth(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer ok.Close()
	if err := checkHealth(ok.URL + "/healthz"); err != nil {
		t.Errorf("200: %v", err)
	}

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer bad.Close()
	if err := checkHealth(bad.URL + "/healthz"); err == nil {
		t.Error("500: no error")
	}

	// Niko ne sluša: greška veze, kao kad aplikacija ne radi.
	url := bad.URL
	bad.Close()
	if err := checkHealth(url + "/healthz"); err == nil {
		t.Error("closed server: no error")
	}
}
