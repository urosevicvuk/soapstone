// Paket message opisuje jedinu stvar sa kojom aplikacija radi: poruku sa
// autorom i tekstom, i pravila koja važe za nju.
package message

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// Granice dužine, u znakovima (rune), a ne u bajtovima: „š“ je jedan znak,
// iako u UTF-8 zauzima dva bajta.
const (
	MaxAuthor = 40
	MaxText   = 280
)

// Message je ono što API vraća. Imena u `json:"…"` su imena polja u JSON-u.
type Message struct {
	ID        int64     `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

// Input je ono što klijent šalje kad pravi poruku: samo autor i tekst.
// ID i vreme dodeljuje skladište.
type Input struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

// Greške validacije. Njihov tekst ide klijentu u polju "error".
var (
	ErrAuthorRequired = errors.New("author is required")
	ErrAuthorTooLong  = errors.New("author must be at most 40 characters")
	ErrTextRequired   = errors.New("text is required")
	ErrTextTooLong    = errors.New("text must be at most 280 characters")
)

// Normalize skida razmake sa početka i kraja oba polja. Validacija i upis
// rade nad ovim oblikom, pa "  Ana  " i "Ana" znače isto.
func (in Input) Normalize() Input {
	return Input{
		Author: strings.TrimSpace(in.Author),
		Text:   strings.TrimSpace(in.Text),
	}
}

// Validate proverava normalizovanu poruku i vraća prvu grešku koju nađe,
// ili nil ako je poruka ispravna.
func (in Input) Validate() error {
	in = in.Normalize()
	switch author, text := utf8.RuneCountInString(in.Author), utf8.RuneCountInString(in.Text); {
	case author == 0:
		return ErrAuthorRequired
	case author > MaxAuthor:
		return ErrAuthorTooLong
	case text == 0:
		return ErrTextRequired
	case text > MaxText:
		return ErrTextTooLong
	}
	return nil
}
