.PHONY: build

BUILD_ARCH ?= $(shell uname)
BUILD_COMMIT ?= $(shell git rev-list -1 HEAD)
BUILD_TIME ?= $(shell date)
BUILD_VERSION ?= $(shell cat BUILD_VERSION)

BINARY_NAME := tldr
PLATFORMS := windows linux darwin
os = $(word 1, $@)

build:
	go build \
		-ldflags "-X \"main.buildTime=$(BUILD_TIME)\" -X \"main.buildVersion=$(BUILD_VERSION)\" -X main.buildCommit=$(BUILD_COMMIT) -X \"main.buildArch=$(BUILD_ARCH)/amd64\"" \
		-o bin/$(BINARY_NAME) cli/tldr/*.go

clean:
	rm bin/$(BINARY_NAME)*

run:
	go run cli/tldr/*.go

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build \
		-ldflags "-X \"main.buildTime=$(BUILD_TIME)\" -X \"main.buildVersion=$(BUILD_VERSION)\" -X \"main.buildCommit=$(BUILD_COMMIT)\" -X \"main.buildArch=$(os)/amd64\"" \
		-o bin/$(BINARY_NAME)-$(VERSION)-$(os)-amd64 cli/tldr/*.go

.PHONY: release
release: windows linux darwin

