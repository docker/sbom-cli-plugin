BIN = docker-sbom
REPO = docker-sbom-cli-plugin

TEMP_DIR = ./.tmp
SNAPSHOT_DIR=./snapshot
DIST_DIR=./dist
LINT_CMD = $(TEMP_DIR)/golangci-lint run --tests=false --timeout=2m --config .golangci.yaml
RELEASE_CMD=$(TEMP_DIR)/goreleaser release --rm-dist
SNAPSHOT_CMD=$(RELEASE_CMD) --skip-publish --rm-dist --snapshot

BOLD := $(shell tput -T linux bold)
PURPLE := $(shell tput -T linux setaf 5)
GREEN := $(shell tput -T linux setaf 2)
CYAN := $(shell tput -T linux setaf 6)
RED := $(shell tput -T linux setaf 1)
RESET := $(shell tput -T linux sgr0)
TITLE := $(BOLD)$(PURPLE)
SUCCESS := $(BOLD)$(GREEN)

## Variable assertions

ifndef TEMP_DIR
	$(error TEMP_DIR is not set)
endif

ifndef SNAPSHOT_DIR
	$(error SNAPSHOT_DIR is not set)
endif

define title
    @printf '$(TITLE)$(1)$(RESET)\n'
endef

## Tasks

.PHONY: all
all: clean static-analysis test ## Run all linux-based checks (linting, license check, unit, integration, and linux acceptance tests)
	@printf '$(SUCCESS)All checks pass!$(RESET)\n'

.PHONY: test
test: unit install-test ## Run all tests

.PHONY: bootstrap-tools
bootstrap-tools: $(TEMP_DIR)

$(TEMP_DIR):
	$(call title,Bootstrapping tools)
	mkdir -p $(TEMP_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TEMP_DIR)/ v1.42.1
	curl -sSfL https://raw.githubusercontent.com/wagoodman/go-bouncer/master/bouncer.sh | sh -s -- -b $(TEMP_DIR)/ v0.3.0
	.github/scripts/goreleaser-install.sh -b $(TEMP_DIR)/ v1.3.1

.PHONY: bootstrap-go
bootstrap-go:
	go mod download

.PHONY: bootstrap
bootstrap: bootstrap-go bootstrap-tools ## Download and install all go dependencies (+ prep tooling in the ./tmp dir)
	$(call title,Bootstrapping dependencies)

.PHONY: static-analysis
static-analysis: lint check-go-mod-tidy check-licenses

.PHONY: lint
lint: ## Run gofmt + golangci lint checks
	$(call title,Running linters)
	# ensure there are no go fmt differences
	@printf "files with gofmt issues: [$(shell gofmt -l -s .)]\n"
	@test -z "$(shell gofmt -l -s .)"

	# run all golangci-lint rules
	$(LINT_CMD)

	# go tooling does not play well with certain filename characters, ensure the common cases don't result in future "go get" failures
	$(eval MALFORMED_FILENAMES := $(shell find . | grep -e ':'))
	@bash -c "[[ '$(MALFORMED_FILENAMES)' == '' ]] || (printf '\nfound unsupported filename characters:\n$(MALFORMED_FILENAMES)\n\n' && false)"

.PHONY: lint-fix
lint-fix: ## Auto-format all source code + run golangci lint fixers
	$(call title,Running lint fixers)
	gofmt -w -s .
	$(LINT_CMD) --fix
	go mod tidy

.PHONY: check-licenses
check-licenses: ## Ensure transitive dependencies are compliant with the current license policy
	$(TEMP_DIR)/bouncer check

check-go-mod-tidy:
	@ .github/scripts/go-mod-tidy-check.sh && echo "go.mod and go.sum are tidy!"

.PHONY: unit
unit:  ## Run unit tests
	$(call title,Running unit tests)
	go test $(shell go list ./... | grep -v anchore/$(REPO)/test)

# note: this is used by CI to determine if the install test fixture cache (docker image tars) should be busted
install-fingerprint:
	cd test/install && \
		make cache.fingerprint

install-test:
	cd test/install && \
		make

install-test-cache-save:
	cd test/install && \
		make save

install-test-cache-load:
	cd test/install && \
		make load

install-test-ci-mac:
	cd test/install && \
		make ci-test-mac

$(SNAPSHOT_DIR): $(TEMP_DIR) ## Build snapshot release binaries and packages
	$(call title,Building snapshot artifacts)
	# create a config with the dist dir overridden
	echo "dist: $(SNAPSHOT_DIR)" > $(TEMP_DIR)/goreleaser.yaml
	cat .goreleaser.yaml >> $(TEMP_DIR)/goreleaser.yaml

	$(SNAPSHOT_CMD) --config $(TEMP_DIR)/goreleaser.yaml

.PHONY: release
release: clean-dist ## Build and publish final binaries and packages
	$(call title,Publishing release artifacts)
	$(RELEASE_CMD)

.PHONY: clean
clean: clean-dist clean-snapshot clean-changelog ## Remove previous builds, result reports, and caches

.PHONY: clean-snapshot
clean-snapshot:
	rm -rf $(SNAPSHOT_DIR) $(TEMP_DIR)/goreleaser.yaml

.PHONY: clean-dist
clean-dist: clean-changelog
	rm -rf $(DIST_DIR) $(TEMP_DIR)/goreleaser.yaml

.PHONY: clean-changelog
clean-changelog:
	rm -f CHANGELOG.md

.PHONY: clean-tmp
clean-tmp:
	rm -rf $(TEMP_DIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BOLD)$(CYAN)%-25s$(RESET)%s\n", $$1, $$2}'
