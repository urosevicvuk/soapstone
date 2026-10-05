package store

import (
	"context"
	"sync"
	"time"

	"github.com/urosevicvuk/soapstone/internal/message"
)

// Memory čuva poruke u memoriji procesa. Brzo je i ne treba mu ništa, ali
// sve nestaje kad se proces ugasi, i svaka kopija aplikacije ima svoje poruke.
type Memory struct {
	// mu štiti polja ispod: HTTP server svaki zahtev obrađuje u posebnoj
	// gorutini, pa više zahteva može istovremeno da čita i menja poruke.
	mu       sync.RWMutex
	messages []message.Message // od najstarije ka najnovijoj
	nextID   int64

	now func() time.Time // sat, zamenljiv u testovima
}

// NewMemory pravi prazno skladište; prva poruka dobija id 1.
func NewMemory() *Memory {
	return &Memory{nextID: 1, now: time.Now}
}

func (m *Memory) Kind() string { return "memory" }

func (m *Memory) Ready(context.Context) error { return nil }

func (m *Memory) List(_ context.Context, limit int) ([]message.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Kopija u obrnutom redu: najnovija prva. Vraćamo kopiju da pozivalac
	// ne bi mogao da menja naš niz bez mutex-a.
	out := make([]message.Message, 0, min(limit, len(m.messages)))
	for i := len(m.messages) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, m.messages[i])
	}
	return out, nil
}

func (m *Memory) Create(_ context.Context, in message.Input) (message.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	msg := message.Message{
		ID:        m.nextID,
		Author:    in.Author,
		Text:      in.Text,
		CreatedAt: m.now().UTC().Truncate(time.Second),
	}
	m.nextID++
	m.messages = append(m.messages, msg)
	return msg, nil
}

func (m *Memory) Get(_ context.Context, id int64) (message.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, msg := range m.messages {
		if msg.ID == id {
			return msg, nil
		}
	}
	return message.Message{}, ErrNotFound
}

func (m *Memory) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, msg := range m.messages {
		if msg.ID == id {
			m.messages = append(m.messages[:i], m.messages[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
