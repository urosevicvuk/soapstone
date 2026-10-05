package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/urosevicvuk/soapstone/internal/message"
	"github.com/urosevicvuk/soapstone/internal/store"
)

// newTestServer pravi server nad praznim memorijskim skladištem, sa
// logovima koji se bacaju, da ne zatrpavaju izlaz testova.
func newTestServer() *Server {
	return New(store.NewMemory(), slog.New(slog.DiscardHandler), "tomato")
}

// do šalje zahtev serveru i vraća odgovor, bez prave mreže.
func do(s *Server, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func TestCRUD(t *testing.T) {
	s := newTestServer()

	// Prva poruka dobija id 1, sa 201 i Location zaglavljem.
	rec := do(s, "POST", "/api/messages", `{"author":"  Solaire ","text":"Praise the Sun"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, body %s", rec.Code, rec.Body)
	}
	var msg message.Message
	if err := json.Unmarshal(rec.Body.Bytes(), &msg); err != nil {
		t.Fatal(err)
	}
	if msg.ID != 1 || msg.Author != "Solaire" || msg.Text != "Praise the Sun" || msg.CreatedAt.IsZero() {
		t.Errorf("created message = %+v", msg)
	}
	if loc := rec.Header().Get("Location"); loc != "/api/messages/1" {
		t.Errorf("Location = %q", loc)
	}

	do(s, "POST", "/api/messages", `{"author":"Ana","text":"Try jumping"}`)

	// Lista ide od najnovije.
	rec = do(s, "GET", "/api/messages", "")
	var list []message.Message
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 || len(list) != 2 || list[0].ID != 2 || list[1].ID != 1 {
		t.Errorf("GET list = %d %+v", rec.Code, list)
	}

	if rec := do(s, "GET", "/api/messages?limit=1", ""); !strings.Contains(rec.Body.String(), "Try jumping") ||
		strings.Contains(rec.Body.String(), "Praise") {
		t.Errorf("limit=1 body = %s", rec.Body)
	}

	if rec := do(s, "GET", "/api/messages/1", ""); rec.Code != 200 {
		t.Errorf("GET 1 status = %d", rec.Code)
	}

	// Brisanje: 204 bez tela, pa 404 i za GET i za ponovni DELETE.
	rec = do(s, "DELETE", "/api/messages/1", "")
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Errorf("DELETE = %d, body %q", rec.Code, rec.Body)
	}
	if rec := do(s, "GET", "/api/messages/1", ""); rec.Code != 404 {
		t.Errorf("GET deleted = %d", rec.Code)
	}
	if rec := do(s, "DELETE", "/api/messages/1", ""); rec.Code != 404 {
		t.Errorf("second DELETE = %d", rec.Code)
	}
}

func TestEmptyListIsArray(t *testing.T) {
	rec := do(newTestServer(), "GET", "/api/messages", "")
	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("empty list = %s, want []", body)
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name, method, target, body string
		status                     int
	}{
		{"invalid JSON", "POST", "/api/messages", `{"author":`, 400},
		{"empty text", "POST", "/api/messages", `{"author":"Ana","text":""}`, 400},
		{"text only spaces", "POST", "/api/messages", `{"author":"Ana","text":"   "}`, 400},
		{"missing author", "POST", "/api/messages", `{"text":"hello"}`, 400},
		{"text too long", "POST", "/api/messages", `{"author":"Ana","text":"` + strings.Repeat("x", 281) + `"}`, 400},
		{"limit zero", "GET", "/api/messages?limit=0", "", 400},
		{"limit too big", "GET", "/api/messages?limit=201", "", 400},
		{"limit not a number", "GET", "/api/messages?limit=all", "", 400},
		{"unknown id", "GET", "/api/messages/99", "", 404},
		{"id not a number", "GET", "/api/messages/abc", "", 404},
		{"delete unknown", "DELETE", "/api/messages/99", "", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(newTestServer(), tt.method, tt.target, tt.body)
			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
			// Svaka greška API-ja je JSON sa poljem "error".
			var e map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &e); err != nil || e["error"] == "" {
				t.Errorf("error body = %s", rec.Body)
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	rec := do(newTestServer(), "GET", "/healthz", "")
	if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Errorf("healthz = %d %q", rec.Code, rec.Body)
	}
}

func TestIndex(t *testing.T) {
	s := newTestServer()
	do(s, "POST", "/api/messages", `{"author":"Ana","text":"<b>Try jumping</b>"}`)

	rec := do(s, "GET", "/", "")
	page := rec.Body.String()

	if rec.Code != 200 || !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("GET / = %d %s", rec.Code, rec.Header().Get("Content-Type"))
	}
	// Boja i skladište se vide, a HTML iz poruke je escape-ovan.
	for _, want := range []string{"background: tomato", "memory", "&lt;b&gt;Try jumping&lt;/b&gt;"} {
		if !strings.Contains(page, want) {
			t.Errorf("page does not contain %q", want)
		}
	}

	// Samo "/" je stranica; ostale nepoznate putanje su 404.
	if rec := do(s, "GET", "/missing", ""); rec.Code != 404 {
		t.Errorf("GET /missing = %d", rec.Code)
	}
}

func TestRequestID(t *testing.T) {
	s := newTestServer()

	// Bez zaglavlja, server ga pravi.
	if id := do(s, "GET", "/healthz", "").Header().Get("X-Request-ID"); id == "" {
		t.Error("no generated X-Request-ID")
	}

	// Sa zaglavljem, server vraća isti.
	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("X-Request-ID", "abc123")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if id := rec.Header().Get("X-Request-ID"); id != "abc123" {
		t.Errorf("X-Request-ID = %q, want abc123", id)
	}
}
