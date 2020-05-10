SHELL=/bin/bash -o pipefail
export GO111MODULE=on

GO ?= go
TEST_FLAGS ?= -v -race

.PHONY: all
all: test build

.PHONY: test
test:
	$(GO) test ${TEST_FLAGS} ./...

.PHONY: coverage
coverage:
	$(GO) test -coverprofile cp.out
	$(GO) tool cover -html=cp.out

.PHONY: build
build:
	$(GO) build -o build/cloudkitchen

.PHONY: clean
clean:
	-rm -rf build cp.out
