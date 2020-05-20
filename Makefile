all: cli

cli:
	@go build .

test:
	@go test -coverprofile=profile.out ./...

.PHONY: cli
