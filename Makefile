.PHONY: dev test

build:
	go build -o bin/server cmd/server/main.go
serve:
	go run cmd/server/main.go
serve-watch:
	air
lint:
	golangci-lint run ./...
test:
	go test -v ./...
test-watch:
	gotestsum --watch
