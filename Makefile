DOCKER = $(shell which docker)
BUILDDIR ?= $(CURDIR)/build

PACKAGES_E2E=$(shell go list ./... | grep '/itest')

ldflags := $(LDFLAGS)
build_tags := $(BUILD_TAGS)
build_args := $(BUILD_ARGS)

ifeq ($(VERBOSE),true)
	build_args += -v
endif

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static" -v
endif

BUILD_TARGETS := build install
BUILD_FLAGS := --tags "$(build_tags)" --ldflags '$(ldflags)'

all: build install

build: BUILD_ARGS := $(build_args) -o $(BUILDDIR)

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-docker:
	$(DOCKER) build --tag scalarorg/protocol-signer -f Dockerfile \
		$(shell git rev-parse --show-toplevel)

.PHONY: build build-docker install tests

test:
	go test ./...

test-e2e:
	go test -mod=readonly -timeout=25m -v $(PACKAGES_E2E) -count=1 --tags=e2e











# Update this path to where bitcoin-vault is located in your system
BITCOIN_VAULT_PATH := $(shell go list -f '{{.Dir}}' github.com/scalarorg/bitcoin-vault/ffi/go-psbt)

.PHONY: build run

# For MacOS
build-darwin:
	CGO_LDFLAGS="-L$(BITCOIN_VAULT_PATH)/lib/darwin -lbitcoin_vault_ffi" \
	CGO_CFLAGS="-I$(BITCOIN_VAULT_PATH)/lib/darwin" \
	go build -o bin/main ./go-ffi/main.go

# For Linux
build-linux:
	CGO_LDFLAGS="-L$(BITCOIN_VAULT_PATH)/lib/linux -lbitcoin_vault_ffi" \
	CGO_CFLAGS="-I$(BITCOIN_VAULT_PATH)/lib/linux" \
	go build -o bin/main ./go-ffi/main.go

# Run the application (MacOS)
run-darwin: build-darwin
	DYLD_LIBRARY_PATH=$(BITCOIN_VAULT_PATH)/lib/darwin ./bin/main

# Run the application (Linux)
run-linux: build-linux
	LD_LIBRARY_PATH=$(BITCOIN_VAULT_PATH)/lib/linux ./bin/main