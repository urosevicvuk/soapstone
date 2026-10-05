package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/urosevicvuk/soapstone/internal/message"
	"github.com/urosevicvuk/soapstone/internal/store"
)

const (
	defaultLimit = 50
	maxLimit     = 200
	maxBody      = 64 << 10 // 64 KiB je mnogo više nego što ispravna poruka traži
)

// handleList: GET /api/messages?limit=50 vraća poruke od najnovije.
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	limit := defaultLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > maxLimit {
			s.writeError(w, http.StatusBadRequest, "limit must be an integer from 1 to 200")
			return
		}
		limit = n
	}

	list, err := s.store.List(r.Context(), limit)
	if err != nil {
		s.internalError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, list)
}

// handleCreate: POST /api/messages sa telom {"author":"…","text":"…"}.
func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var in message.Input
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody)).Decode(&in); err != nil {
		s.writeError(w, http.StatusBadRequest, `body must be JSON like {"author": "…", "text": "…"}`)
		return
	}
	if err := in.Validate(); err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := s.store.Create(r.Context(), in.Normalize())
	if err != nil {
		s.internalError(w, err)
		return
	}

	// 201 Created, a Location kaže gde se nova poruka čita.
	w.Header().Set("Location", "/api/messages/"+strconv.FormatInt(msg.ID, 10))
	s.writeJSON(w, http.StatusCreated, msg)
}

// handleGet: GET /api/messages/{id}.
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id, ok := s.pathID(w, r)
	if !ok {
		return
	}
	msg, err := s.store.Get(r.Context(), id)
	if err != nil {
		s.storeError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, msg)
}

// handleDelete: DELETE /api/messages/{id} vraća 204 No Content, bez tela.
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := s.pathID(w, r)
	if !ok {
		return
	}
	if err := s.store.Delete(r.Context(), id); err != nil {
		s.storeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// pathID čita {id} iz putanje. Id koji nije ceo broj ne može da postoji,
// pa je odgovor 404, isto kao za nepostojeći broj.
func (s *Server) pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.writeError(w, http.StatusNotFound, store.ErrNotFound.Error())
		return 0, false
	}
	return id, true
}

// storeError pretvara grešku skladišta u odgovor: ErrNotFound je 404, a
// sve ostalo je greška servera.
func (s *Server) storeError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		s.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	s.internalError(w, err)
}

// internalError loguje pravu grešku, a klijentu šalje samo da je do servera.
func (s *Server) internalError(w http.ResponseWriter, err error) {
	s.log.Error("internal error", "error", err)
	s.writeError(w, http.StatusInternalServerError, "internal server error")
}
