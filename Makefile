.PHONY: build

BINARY_NAME := tldr
VERSION ?= $(shell cat VERSION)
TODAY ?= $(shell date)
GIT_COMMIT ?= $(shell git rev-list -1 HEAD)
ARCH ?= $(shell uname)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

build:
	go build \
		-ldflags "-X main.buildDate=$(TODAY) -X main.buildVersion=$(VERSION) -X main.buildCommit=$(GIT_COMMIT) -X main.buildArch=$(ARCH)/amd64" \
		-o bin/$(BINARY_NAME) cli/tldr/*.go

clean:
	rm bin/$(BINARY_NAME)*

run:
	go run cli/tldr/*.go

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build \
		-ldflags "-X main.buildDate=$(TODAY) -X main.buildVersion=$(VERSION) -X main.buildCommit=$(GIT_COMMIT) -X main.buildArch=$(os)/amd64" \
		-o bin/$(BINARY_NAME)-$(VERSION)-$(os)-amd64 cli/tldr/*.go

.PHONY: release
release: windows linux darwin

