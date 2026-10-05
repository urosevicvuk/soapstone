package store

import (
	"context"
	"errors"
	"testing"

	"github.com/urosevicvuk/soapstone/internal/message"
)

func TestMemory(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()

	for _, text := range []string{"first", "second", "third"} {
		if _, err := m.Create(ctx, message.Input{Author: "Ana", Text: text}); err != nil {
			t.Fatal(err)
		}
	}

	// Lista ide od najnovije, a limit je skraćuje.
	list, _ := m.List(ctx, 2)
	if len(list) != 2 || list[0].Text != "third" || list[1].Text != "second" {
		t.Fatalf("List(2) = %+v", list)
	}

	// Id-jevi kreću od 1.
	msg, err := m.Get(ctx, 1)
	if err != nil || msg.Text != "first" {
		t.Fatalf("Get(1) = %+v, %v", msg, err)
	}

	if err := m.Delete(ctx, 1); err != nil {
		t.Fatalf("Delete(1) = %v", err)
	}
	if _, err := m.Get(ctx, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(1) after delete = %v, want ErrNotFound", err)
	}
	if err := m.Delete(ctx, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second Delete(1) = %v, want ErrNotFound", err)
	}

	// Obrisan id se ne koristi ponovo.
	msg, _ = m.Create(ctx, message.Input{Author: "Ana", Text: "fourth"})
	if msg.ID != 4 {
		t.Fatalf("id after delete = %d, want 4", msg.ID)
	}
}
