package buildinfo

import (
	"strings"
	"testing"
)

func TestGetDefaults(t *testing.T) {
	i := Get()
	if i.Version != "dev" || i.Commit != "unknown" {
		t.Errorf("Get() = %+v, want dev and unknown without -ldflags", i)
	}
	if !strings.HasPrefix(i.GoVersion, "go") {
		t.Errorf("GoVersion = %q, want go1.x", i.GoVersion)
	}
	if got := i.String(); got != "dev (unknown)" {
		t.Errorf("String() = %q", got)
	}
}
