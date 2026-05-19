BIN = docker-sbom
REPO = sbom-cli-plugin

TEMP_DIR = ./.tmp
DIST_DIR=./dist
SNAPSHOT_DIR=./snapshot
RESULTS_DIR = test/results
COVER_REPORT = $(RESULTS_DIR)/unit-coverage-details.txt
COVER_TOTAL = $(RESULTS_DIR)/unit-coverage-summary.txt
LINT_CMD = $(TEMP_DIR)/golangci-lint run --tests=false --timeout=2m --config .golangci.yaml
GOIMPORTS_CMD = $(TEMP_DIR)/gosimports -local github.com/anchore
RELEASE_CMD=$(TEMP_DIR)/goreleaser release --rm-dist
RELEASE_FLAGS=""
SNAPSHOT_CMD=$(RELEASE_CMD) --skip-publish --rm-dist --snapshot
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
SNAPSHOT_BIN=$(shell realpath $(shell pwd)/$(SNAPSHOT_DIR)/$(REPO)_$(OS)_amd64/$(BIN))

BOLD := $(shell tput -T linux bold)
PURPLE := $(shell tput -T linux setaf 5)
GREEN := $(shell tput -T linux setaf 2)
CYAN := $(shell tput -T linux setaf 6)
RED := $(shell tput -T linux setaf 1)
RESET := $(shell tput -T linux sgr0)
TITLE := $(BOLD)$(PURPLE)
SUCCESS := $(BOLD)$(GREEN)

## change these values manually if you'd like to bust the cache in CI for select test fixtures
CLI_CACHE_BUSTER = d12f51e6c910590b485b

## Variable assertions

ifndef RESULTS_DIR
	$(error RESULTS_DIR is not set)
endif

ifndef TEMP_DIR
	$(error TEMP_DIR is not set)
endif

ifndef SNAPSHOT_DIR
	$(error SNAPSHOT_DIR is not set)
endif

define title
    @printf '$(TITLE)$(1)$(RESET)\n'
endef

define safe_rm_rf
	bash -c 'test -z "$(1)" && false || rm -rf $(1)'
endef

define safe_rm_rf_children
	bash -c 'test -z "$(1)" && false || rm -rf $(1)/*'
endef

## Tasks

.PHONY: all
all: clean-snapshot static-analysis $(SNAPSHOT_DIR) test ## Run all linux-based checks (linting, license check, unit, integration, and linux acceptance tests)
	@printf '$(SUCCESS)All checks pass!$(RESET)\n'

.PHONY: test
test: unit install-test cli ## Run all tests

$(RESULTS_DIR):
	mkdir -p $(RESULTS_DIR)

.PHONY: bootstrap-tools
bootstrap-tools:
	$(call title,Bootstrapping tools)
	mkdir -p $(TEMP_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TEMP_DIR)/ v1.45.0
	curl -sSfL https://raw.githubusercontent.com/wagoodman/go-bouncer/master/bouncer.sh | sh -s -- -b $(TEMP_DIR)/ v0.3.0
	curl -sSfL https://raw.githubusercontent.com/anchore/chronicle/main/install.sh | sh -s -- -b $(TEMP_DIR)/ v0.4.1
	.github/scripts/goreleaser-install.sh -b $(TEMP_DIR)/ v1.5.0
	# the only difference between goimports and gosimports is that gosimports removes extra whitespace between import blocks (see https://github.com/golang/go/issues/20818)
	GOBIN="$(shell realpath $(TEMP_DIR))" go install github.com/rinchsan/gosimports/cmd/gosimports@v0.1.5

.PHONY: bootstrap-go
bootstrap-go:
	go mod download

.PHONY: bootstrap
bootstrap: $(RESULTS_DIR) bootstrap-go bootstrap-tools ## Download and install all go dependencies (+ prep tooling in the ./tmp dir)
	$(call title,Bootstrapping go dependencies)

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
	@[ -z "$(shell $(GOIMPORTS_CMD) -d .)" ] && echo "goimports clean" || (echo "goimports needs to be fixed" && false)

	# go tooling does not play well with certain filename characters, ensure the common cases don't result in future "go get" failures
	$(eval MALFORMED_FILENAMES := $(shell find . | grep -e ':'))
	@bash -c "[[ '$(MALFORMED_FILENAMES)' == '' ]] || (printf '\nfound unsupported filename characters:\n$(MALFORMED_FILENAMES)\n\n' && false)"

.PHONY: lint-fix
lint-fix: ## Auto-format all source code + run golangci lint fixers
	$(call title,Running lint fixers)
	gofmt -w -s .
	$(GOIMPORTS_CMD) -w .
	$(LINT_CMD) --fix
	go mod tidy

.PHONY: check-licenses
check-licenses: ## Ensure transitive dependencies are compliant with the current license policy
	$(TEMP_DIR)/bouncer check

check-go-mod-tidy:
	@ .github/scripts/go-mod-tidy-check.sh && echo "go.mod and go.sum are tidy!"

.PHONY: unit
unit: $(RESULTS_DIR)  ## Run unit tests
	$(call title,Running unit tests)
	go test  -coverprofile $(COVER_REPORT) $(shell go list ./... | grep -v docker/sbom-cli-plugin/test)
	@go tool cover -func $(COVER_REPORT) | grep total |  awk '{print substr($$3, 1, length($$3)-1)}' > $(COVER_TOTAL)
	@echo "Coverage: $$(cat $(COVER_TOTAL))"

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

# note: this is used by CI to determine if the integration test fixture cache (docker image tars) should be busted
cli-fingerprint:
	$(call title,CLI test fixture fingerprint)
	find test/cli/test-fixtures/image-* -type f -exec md5sum {} + | awk '{print $1}' | sort | md5sum | tee test/cli/test-fixtures/cache.fingerprint && echo "$(CLI_CACHE_BUSTER)" >> test/cli/test-fixtures/cache.fingerprint

.PHONY: cli
cli:
	$(call title,Building snapshot artifacts)
	echo "dist: $(SNAPSHOT_DIR)" > $(TEMP_DIR)/goreleaser.yaml
	cat .goreleaser.yaml >> $(TEMP_DIR)/goreleaser.yaml
	$(SNAPSHOT_CMD) --config $(TEMP_DIR)/goreleaser.yaml

.PHONY: install-snapshot
install-snapshot:
	cp $(SNAPSHOT_BIN) ~/.docker/cli-plugins/

.PHONY: install
install: cli
	cp $(SNAPSHOT_BIN) ~/.docker/cli-plugins/

.PHONY: changelog
changelog: clean-changelog CHANGELOG.md
	@docker run -it --rm \
		-v $(shell pwd)/CHANGELOG.md:/CHANGELOG.md \
		rawkode/mdv \
			-t 748.5989 \
			/CHANGELOG.md

CHANGELOG.md:
	$(TEMP_DIR)/chronicle -vv > CHANGELOG.md

.PHONY: validate-syft-release-version
validate-syft-release-version:
	@./.github/scripts/syft-released-version-check.sh

.PHONY: release
release: clean-dist CHANGELOG.md
	$(call title,Publishing release artifacts)
	bash -c "$(RELEASE_CMD) $(RELEASE_FLAGS) --release-notes <(cat CHANGELOG.md)"

.PHONY: clean
clean: clean-dist clean-snapshot  ## Remove previous builds, result reports, and test cache
	$(call safe_rm_rf_children,$(RESULTS_DIR))

.PHONY: clean-snapshot
clean-snapshot:
	$(call safe_rm_rf,$(SNAPSHOT_DIR))
	rm -f $(TEMP_DIR)/goreleaser.yaml

.PHONY: clean-dist
clean-dist: clean-changelog
	$(call safe_rm_rf,$(DIST_DIR))
	rm -f $(TEMP_DIR)/goreleaser.yaml

.PHONY: clean-changelog
clean-changelog:
	rm -f CHANGELOG.md


.PHONY: clean-tmp
clean-tmp:
	rm -rf $(TEMP_DIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BOLD)$(CYAN)%-25s$(RESET)%s\n", $$1, $$2}'
