.PHONY: build test lint docker

BINARY=indexer

build:
	go build -v -o $(BINARY) ./cmd/indexer

lint:
	golangci-lint run

test:
	go test ./... -v

docker:
	docker build -t content-indexer:latest .