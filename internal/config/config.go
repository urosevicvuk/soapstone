// Paket config čita podešavanja iz env promenljivih. Aplikacija ne čita
// nikakav konfiguracioni fajl: sve se zadaje spolja, isto lokalno, u
// kontejneru i u klasteru.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

// Config su sva podešavanja aplikacije, već proverena.
type Config struct {
	Port      int
	AppColor  string
	LogLevel  slog.Level
	LogFormat string // "json" ili "text"
}

// Load čita env promenljive. Za neispravnu vrednost vraća grešku koja kaže
// koja promenljiva i zašto, umesto da aplikacija tiho krene sa pogrešnim
// podešavanjem.
func Load() (Config, error) {
	return load(os.Getenv)
}

// load prima funkciju za čitanje promenljivih, da bi testovi mogli da
// proslede svoje vrednosti bez diranja pravog okruženja.
func load(getenv func(string) string) (Config, error) {
	get := func(key, def string) string {
		if v := getenv(key); v != "" {
			return v
		}
		return def
	}

	var c Config

	port, err := strconv.Atoi(get("PORT", "8080"))
	if err != nil || port < 1 || port > 65535 {
		return c, fmt.Errorf("PORT must be a number from 1 to 65535, got %q", getenv("PORT"))
	}
	c.Port = port

	c.AppColor = get("APP_COLOR", "steelblue")

	switch lvl := get("LOG_LEVEL", "info"); lvl {
	case "debug":
		c.LogLevel = slog.LevelDebug
	case "info":
		c.LogLevel = slog.LevelInfo
	case "warn":
		c.LogLevel = slog.LevelWarn
	case "error":
		c.LogLevel = slog.LevelError
	default:
		return c, fmt.Errorf("LOG_LEVEL must be debug, info, warn or error, got %q", lvl)
	}

	switch f := get("LOG_FORMAT", "json"); f {
	case "json", "text":
		c.LogFormat = f
	default:
		return c, fmt.Errorf("LOG_FORMAT must be json or text, got %q", f)
	}

	return c, nil
}
