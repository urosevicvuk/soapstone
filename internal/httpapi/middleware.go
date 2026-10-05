package httpapi

import (
	"crypto/rand"
	"net/http"
	"time"
)

// statusRecorder pamti status koji je handler poslao, jer ga
// http.ResponseWriter posle ne otkriva.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Unwrap daje pristup originalnom ResponseWriter-u (http.ResponseController).
func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// logRequests ispisuje jednu liniju loga po zahtevu, posle odgovora, i
// svakom zahtevu daje request id.
//
// Request id je oznaka po kojoj se jedan zahtev prati kroz logove. Ako ga
// klijent ili proxy ispred nas već pošalje u X-Request-ID, koristimo njega;
// inače ga pravimo. U oba slučaja vraća se u istom zaglavlju odgovora.
func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = rand.Text() // nasumičan niz od 26 znakova
		}
		w.Header().Set("X-Request-ID", id)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// r.Pattern je ruta koju je ServeMux izabrao, na primer
		// "GET /api/messages/{id}", a ne pun path sa konkretnim id-jem.
		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}

		s.log.Info("request",
			"method", r.Method,
			"route", route,
			"status", rec.status,
			"duration_ms", float64(time.Since(start).Microseconds())/1000,
			"request_id", id,
		)
	})
}
