package message

import (
	"errors"
	"strings"
	"testing"
)

// Table-driven test: svaki red tabele je jedan slučaj, a petlja ih sve
// pokreće kao podtestove. Nov slučaj je samo nov red.
func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		input Input
		want  error
	}{
		{"valid", Input{Author: "Solaire", Text: "Praise the Sun"}, nil},
		{"surrounding spaces ignored", Input{Author: "  Solaire  ", Text: "  Try jumping  "}, nil},

		{"empty author", Input{Author: "", Text: "Try jumping"}, ErrAuthorRequired},
		{"author only spaces", Input{Author: "   ", Text: "Try jumping"}, ErrAuthorRequired},
		{"author exactly 40", Input{Author: strings.Repeat("a", 40), Text: "x"}, nil},
		{"author 41", Input{Author: strings.Repeat("a", 41), Text: "x"}, ErrAuthorTooLong},
		{"author 40 multibyte", Input{Author: strings.Repeat("š", 40), Text: "x"}, nil},

		{"empty text", Input{Author: "Ana", Text: ""}, ErrTextRequired},
		{"text only whitespace", Input{Author: "Ana", Text: " \n\t "}, ErrTextRequired},
		{"text exactly 1", Input{Author: "Ana", Text: "x"}, nil},
		{"text exactly 280", Input{Author: "Ana", Text: strings.Repeat("x", 280)}, nil},
		{"text 281", Input{Author: "Ana", Text: strings.Repeat("x", 281)}, ErrTextTooLong},
		{"text 280 multibyte", Input{Author: "Ana", Text: strings.Repeat("ž", 280)}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Validate()
			if !errors.Is(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	got := Input{Author: " Ana ", Text: "\thello\n"}.Normalize()
	want := Input{Author: "Ana", Text: "hello"}
	if got != want {
		t.Errorf("Normalize() = %+v, want %+v", got, want)
	}
}
