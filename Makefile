all: cli

cli:
	@go build .

test:
	@go test -cover ./...

.PHONY: cli
