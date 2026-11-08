package httpapi

import (
	"net/http"

	"github.com/urosevicvuk/soapstone/internal/buildinfo"
)

// handleVersion vraća verziju i commit iz kojih je binary napravljen, da se
// sa servera može pročitati šta tačno na njemu radi.
func (s *Server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	s.writeJSON(w, http.StatusOK, buildinfo.Get())
}
