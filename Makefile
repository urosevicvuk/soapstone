# Komande za rad sa aplikacijom. `make <cilj>` izvršava recept ispod cilja,
# a `make -n <cilj>` samo ispiše komande, bez izvršavanja.

APP := soapstone

.PHONY: run build test fmt

# Pokreće server na portu 8080 (ili PORT). `run: build` znači da make prvo
# uradi cilj build. Ako postoji .env, njegove vrednosti se izvoze kao env
# promenljive pre pokretanja: `set -a` izvozi sve što se postavi, a `. ./.env`
# učitava fajl u trenutni shell. `exec` zamenjuje shell aplikacijom, pa Ctrl+C
# stiže pravo do nje i ona se gasi uredno.
run: build
	set -a; if [ -f .env ]; then . ./.env; fi; set +a; exec ./bin/$(APP)

# Pravi binarni fajl bin/soapstone.
build:
	go build -o bin/$(APP) ./cmd/$(APP)

# Pokreće sve testove.
test:
	go test ./...

# Formatira sav Go kod.
fmt:
	gofmt -w .
