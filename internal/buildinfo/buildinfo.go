// Paket buildinfo nosi verziju i commit iz kojih je binary napravljen.
//
// Vrednosti se ne pišu u kod, nego se upisuju pri build-u, sa
//
//	go build -ldflags "-X github.com/urosevicvuk/soapstone/internal/buildinfo.Version=v0.2.0 \
//	  -X github.com/urosevicvuk/soapstone/internal/buildinfo.Commit=3f9c2a1"
//
// Bez toga ostaju dev i unknown, pa se odmah vidi da binary nije napravljen
// za objavu.
package buildinfo

import "runtime"

// Version i Commit su promenljive, ne konstante, jer -X ume da menja samo
// promenljive tipa string.
var (
	Version = "dev"
	Commit  = "unknown"
)

// Info je ono što vraćaju `soapstone version` i GET /version.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	GoVersion string `json:"goVersion"`
}

// Get vraća verziju, commit i verziju Go-a kojim je binary napravljen.
func Get() Info {
	return Info{Version: Version, Commit: Commit, GoVersion: runtime.Version()}
}

// String je kratak oblik za terminal: v0.2.0 (3f9c2a1).
func (i Info) String() string {
	return i.Version + " (" + i.Commit + ")"
}
