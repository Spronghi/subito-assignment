.PHONY: dev test

dev:
	air
test:
	gotestsum --watch
