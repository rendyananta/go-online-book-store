build-db:
	CGO_ENABLED=1 go build -tags fts5 -o cmd/bin github.com/rendyananta/example-online-book-store/cmd/db

db-refresh: build-db
	./cmd/bin/db down
	./cmd/bin/db up

build-seeder:
	CGO_ENABLED=1 go build -tags "fts5" -o cmd/bin github.com/rendyananta/example-online-book-store/cmd/seed

db-seed: build-seeder
	./cmd/bin/seed --file cmd/seed/books.csv

build-http:
	AUTH_CIPHER_KEYS=1fYGJsZSQuI0EQEbnCkkMYIh78epX7Tb CGO_ENABLED=1 go build -tags "fts5" -o cmd/bin github.com/rendyananta/example-online-book-store/cmd/http

run-http: build-http
	./cmd/bin/http

test:
	CGO_ENABLED=1 go test ./...
