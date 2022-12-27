ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)
TAG := v1.0.0-beta.3

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

ifeq ($(OSTYPE),Darwin)
    export MallocNanoZone=0
endif

include ./misc/makefile/tools.Makefile

build: test
	@go build ./...

install-deps: gotestsum tparse ## Install Development Dependencies (localy).
deps: $(GOTESTSUM) $(TPARSE) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile   coverage.out
TESTS_ARGS += -test.timeout        60s
TESTS_ARGS += -race
run-tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short

test: run-tests $(TPARSE) ## Run Tests & parse details
	@cat gotestsum.json.out | $(TPARSE) -all -notests


lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

.PHONY: lint lint-prepare clean build unittest

go_mod_tidy:
	go get -u && go mod tidy -go=1.19
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy -go=1.19); \
	done

release:
	TAG=$(TAG) ./scripts/release.sh
	TAG=$(TAG) ./scripts/tag.sh
