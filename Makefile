.PHONY: dev test

serve-watch:
	air
test-watch:
	gotestsum --watch
