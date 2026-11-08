package httpapi

import (
	"encoding/json"
	"testing"

	"github.com/urosevicvuk/soapstone/internal/buildinfo"
)

func TestVersion(t *testing.T) {
	rec := do(newTestServer(), "GET", "/version", "")
	if rec.Code != 200 {
		t.Fatalf("GET /version = %d", rec.Code)
	}
	var got buildinfo.Info
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("GET /version: %v in %q", err, rec.Body)
	}
	if got != buildinfo.Get() {
		t.Errorf("GET /version = %+v, want %+v", got, buildinfo.Get())
	}
}
