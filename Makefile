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
test-no-cache:
	go test -count=1 -v ./...
test-watch:
	gotestsum --watch
