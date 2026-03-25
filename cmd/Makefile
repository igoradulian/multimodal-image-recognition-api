APP=ai_bot
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

build:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(APP) .
build-x86_64:
	$(MAKE) build GOARCH=amd64

.PHONY: build build-x86 build-x86_64
