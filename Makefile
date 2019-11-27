TOOLS_PKG := github.com/dantin/mserver

ARCHIVE_FILE := mserver.tar.gz
VERSION      ?= "1.0.0+git"

PACKAGES := $$(go list ./... | grep -vE 'vendor|cmd')

LDFLAGS += -X "$(TOOLS_PKG)/server.Version=$(shell echo $(VERSION))"
LDFLAGS += -X "$(TOOLS_PKG)/server.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(TOOLS_PKG)/server.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/server.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

GOMOD := -mod=vendor

GOVER_MAJOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\1/")
GOVER_MINOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\2/")
GO111 := $(shell [ $(GOVER_MAJOR) -gt 1  ] || [ $(GOVER_MAJOR) -eq 1  ] && [ $(GOVER_MINOR) -ge 11  ]; echo $$?)
ifeq ($(GO111), 1)
	$(warning "go below 1.11 does not support modules")
GOMOD :=
endif

DEFAULT: build

test:
	@GO111MODULE=on go test --race --cover $(PACKAGES)

.PHONY: build
build:
	@echo "build binary for MacOS"
	@GO111MODULE=on CGO_ENABLED=0 go build $(GOMOD) -ldflags '$(LDFLAGS)'\
		-o bin/media-server cmd/media-server/main.go

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
	@bash $$PWD/scripts/clean_vendor.sh

.PHONY: archive
archive:
	@a=$(ARCHIVE_FILE); \
		tar -zcp \
		--exclude='.git' \
		--exclude='$(ARCHIVE_FILE)' \
		--exclude='bin' \
		-f ../$$a .; \
		mv ../$$a .

.PHONY: clean
clean:
	@echo "clean project"
	@rm -f $(ARCHIVE_FILE)
	@rm -rf bin
