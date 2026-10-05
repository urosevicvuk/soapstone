package httpapi

import (
	"bytes"
	"net/http"

	"github.com/urosevicvuk/soapstone/internal/message"
)

// indexData je ono što šablon web/index.html prikazuje.
type indexData struct {
	Color     string
	Hostname  string
	Store     string
	Messages  []message.Message
	MaxAuthor int
	MaxText   int
}

// handleIndex: GET / prikazuje stranicu sa porukama i formom.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	list, err := s.store.List(r.Context(), defaultLimit)
	if err != nil {
		s.internalError(w, err)
		return
	}

	// Šablon se prvo popunjava u bafer: ako nešto pukne, klijent dobija
	// čist 500, a ne pola stranice.
	var buf bytes.Buffer
	err = s.index.Execute(&buf, indexData{
		Color:     s.color,
		Hostname:  s.hostname,
		Store:     s.store.Kind(),
		Messages:  list,
		MaxAuthor: message.MaxAuthor,
		MaxText:   message.MaxText,
	})
	if err != nil {
		s.internalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		// Zaglavlja su već poslata, pa ostaje samo da se zabeleži; najčešće
		// je klijent prekinuo vezu.
		s.log.Debug("writing page", "error", err)
	}
}
