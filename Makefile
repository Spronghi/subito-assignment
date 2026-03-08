.PHONY: dev test

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
