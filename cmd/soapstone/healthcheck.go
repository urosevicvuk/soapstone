package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/urosevicvuk/soapstone/internal/config"
)

// healthcheck pita /healthz aplikacije koja radi na istoj mašini (u istom
// kontejneru) i vraća grešku ako odgovor nije 200. To je `curl -f` ugrađen u
// binary, za image u kom curl ne postoji: Docker HEALTHCHECK pokreće
// `soapstone healthcheck` i gleda samo exit code, 0 ili 1.
func healthcheck() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return checkHealth("http://localhost:" + strconv.Itoa(cfg.Port) + "/healthz")
}

// checkHealth je sam GET, odvojen da bi ga test pozvao nad svojim serverom.
// Rok je kratak: Docker ionako ponavlja proveru na svakih nekoliko sekundi.
func checkHealth(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned %d", url, resp.StatusCode)
	}
	return nil
}
