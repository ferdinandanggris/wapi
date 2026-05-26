.PHONY: build test lint fmt

build:
	go build ./...

test:
	CGO_ENABLED=1 go test -race ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .
