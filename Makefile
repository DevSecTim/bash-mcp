.PHONY: build run test clean

BINARY_NAME=bash-mcp

build:
	go build -o bin/$(BINARY_NAME) ./cmd/bash-mcp

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -rf bin/

lint:
	golangci-lint run ./...
