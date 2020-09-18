DIST := dist
IMPORT := code.gitea.io/tea
export GO111MODULE=on

GO ?= go
SED_INPLACE := sed -i
SHASUM ?= shasum -a 256

export PATH := $($(GO) env GOPATH)/bin:$(PATH)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		SED_INPLACE := sed -i ''
	endif
endif

GOFILES := $(shell find . -name "*.go" -type f ! -path "./vendor/*" ! -path "*/bindata.go")
GOFMT ?= gofmt -s

MAKE_VERSION := $(shell make -v | head -n 1)

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
	TEA_VERSION ?= $(VERSION)
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
	TEA_VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

TAGS ?=
TAGS_STATIC := osusergo,netgo,static_build,$(TAGS)

LDFLAGS := -X "main.Version=$(TEA_VERSION)" -X "main.Tags=$(TAGS)"
LDFLAGS_STATIC := -X "main.Version=$(TEA_VERSION)" -X "main.Tags=$(TAGS_STATIC)" -linkmode external -extldflags "-fno-PIC -static"

GOFLAGS := -mod=vendor -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' # includes flags for native builds not needed / supported by xgo
XGOFLAGS := -tags '$(TAGS)' -ldflags '$(LDFLAGS) -s -w'         # smaller non-debug binaries
XGOFLAGS_STATIC := -tags '$(TAGS_STATIC)' -ldflags '$(LDFLAGS_STATIC) -s -w' -buildmode=pie
GOFLAGS_STATIC := $(GOFLAGS) $(XGOFLAGS_STATIC)

PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)

ifeq ($(OS), Windows_NT)
	EXECUTABLE := tea.exe
else
	EXECUTABLE := tea
endif

# $(call strip-suffix,filename)
strip-suffix = $(firstword $(subst ., ,$(1)))

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -mod=vendor -i ./...
	rm -rf $(EXECUTABLE) $(DIST)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: vet
vet:
	# Default vet
	$(GO) vet -mod=vendor $(PACKAGES)
	# Custom vet
	$(GO) build -mod=vendor code.gitea.io/gitea-vet
	$(GO) vet -vettool=gitea-vet $(PACKAGES)

.PHONY: lint
lint:
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml -exclude=./vendor/... ./... || exit 1

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error -i unknwon,destory $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w -i unknwon $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: test
test:
	$(GO) test -mod=vendor -tags='sqlite sqlite_unlock_notify' $(PACKAGES)

.PHONY: unit-test-coverage
unit-test-coverage:
	$(GO) test -mod=vendor -tags='sqlite sqlite_unlock_notify' -cover -coverprofile coverage.out $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: vendor
vendor:
	$(GO) mod tidy && $(GO) mod vendor

.PHONY: test-vendor
test-vendor: vendor
	@diff=$$(git diff vendor/); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make vendor' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: check
check: test

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	@echo "building development executable '$@'"
	$(GO) build $(GOFLAGS) -o $@

.PHONY: install
install: $(SOURCES)
	@echo "installing static executable to $(GOPATH)/bin/$(EXECUTABLE)"
	$(GO) install $(GOFLAGS_STATIC)

.PHONY: release
release: release-dirs release-windows release-linux release-darwin release-copy release-compress release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-windows
release-windows:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u src.techknowlogick.com/xgo; \
	fi
	# go modules are turned off due to https://github.com/techknowlogick/xgo/issues/16
	# xgo autodetects --mod=vendor, so no need to set it here.
	GO111MODULE=off xgo -dest $(DIST)/binaries $(XGOFLAGS_STATIC) -targets 'windows/*' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-linux
release-linux:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u src.techknowlogick.com/xgo; \
	fi
	# go modules are turned off due to https://github.com/techknowlogick/xgo/issues/16
	# xgo autodetects --mod=vendor, so no need to set it here.
	GO111MODULE=off xgo -dest $(DIST)/binaries $(XGOFLAGS_STATIC) -targets 'linux/amd64,linux/386,linux/arm-5,linux/arm-6,linux/arm64,linux/mips64le,linux/mips,linux/mipsle' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-darwin
release-darwin:
	@hash xgo > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		cd /tmp && $(GO) get -u src.techknowlogick.com/xgo; \
	fi
	# go modules are turned off due to https://github.com/techknowlogick/xgo/issues/16
	# xgo autodetects --mod=vendor, so no need to set it here.
	# osx requires dynamic linking to system libraries, so don't build a static binary.
	GO111MODULE=off xgo -dest $(DIST)/binaries $(XGOFLAGS) -targets 'darwin/*' -out tea-$(VERSION) .
ifeq ($(CI),drone)
	cp /build/* $(DIST)/binaries
endif

.PHONY: release-copy
release-copy:
	cd $(DIST); for file in `find /build -type f -name "*"`; do cp $${file} ./release/; done;

.PHONY: release-compress
release-compress:
	# go modules are turned off due to https://github.com/techknowlogick/xgo/issues/16
	@hash gxz > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		GO111MODULE=off $(GO) get -u github.com/ulikunitz/xz/cmd/gxz; \
	fi
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && gxz -k -9 $${file}; done;

.PHONY: release-check
release-check:
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;
