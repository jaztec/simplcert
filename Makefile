# Project information
VERSION? := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go build variables
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

CMD := $(GOBASE)/cmd

# Linker flags
LDFLAGS=-v -ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all build clean lint

all: build

build:  ## Build the binary file
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -mod mod $(LDFLAGS) -v -o $(GOBIN)/simplcert $(CMD)/

build-darwin:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=darwin GOARCH=arm64 go build -mod mod $(LDFLAGS) -v -o $(GOBIN)/simplcert-darwin $(CMD)/

test: ## Test the library
	@mkdir -p artifacts/profiles
	go test ./... -bench=. -race -timeout 10000ms -coverprofile artifacts/cover.out
	go tool cover -func=artifacts/cover.out
	go tool cover -html=artifacts/cover.out

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
