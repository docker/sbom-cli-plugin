NAME=docker-sbom-cli-plugin

IMAGE_NAME=$(NAME)-install.sh-env
UBUNTU_IMAGE=$(IMAGE_NAME):ubuntu-20.04
ALPINE_IMAGE=$(IMAGE_NAME):alpine-3.6

ENVS=./environments
DOCKER_RUN=docker run --rm -t -w /project/test/install -v /var/run/docker.sock:/var/run/docker.sock -v $(shell pwd)/../../:/project
UNIT=make unit-local

# acceptance testing is running the current install.sh against the latest release. Note: this could be a problem down
# the line if there are breaking changes made that don't align with the latest release (but will be OK with the next
# release)
CURL_INSTALL_SYFT=curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
WGET_INSTALL_SYFT=wget -O - https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
ACCEPTANCE_CMD=../../install.sh && docker sbom version

# CI cache busting values; change these if you want CI to not use previous stored cache
INSTALL_TEST_CACHE_BUSTER=894d8ca

define title
    @printf '\n≡≡≡[ $(1) ]≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡\n'
endef

.PHONY: test
test: unit acceptance

.PHONY: ci-test-mac
ci-test-mac: unit-local acceptance-local

# note: do not add acceptance-local to this list
acceptance: acceptance-ubuntu-20.04 acceptance-alpine-3.6

unit: unit-ubuntu-20.04

unit-local:
	$(call title,unit tests)
	@for f in $(shell ls *_test.sh); do echo "Running unit test suite '$${f}'"; bash $${f} || exit 1; done

acceptance-local:
	$(CURL_INSTALL_SYFT) && $(ACCEPTANCE_CMD)

save: ubuntu-20.04 alpine-3.6
	@mkdir cache || true
	docker image save -o cache/ubuntu-env.tar $(UBUNTU_IMAGE)
	docker image save -o cache/alpine-env.tar $(ALPINE_IMAGE)

load:
	docker image load -i cache/ubuntu-env.tar
	docker image load -i cache/alpine-env.tar

## UBUNTU #######################################################

acceptance-ubuntu-20.04: ubuntu-20.04
	$(call title,ubuntu:20.04 - acceptance)
	$(DOCKER_RUN) $(UBUNTU_IMAGE) \
		sh -c '$(CURL_INSTALL_SYFT) && $(ACCEPTANCE_CMD)'

unit-ubuntu-20.04: ubuntu-20.04
	$(call title,ubuntu:20.04 - unit)
	$(DOCKER_RUN) $(UBUNTU_IMAGE) \
		$(UNIT)

ubuntu-20.04:
	$(call title,ubuntu:20.04 - build environment)
	docker build -t $(UBUNTU_IMAGE) -f $(ENVS)/Dockerfile-ubuntu-20.04 .

## ALPINE #######################################################

# note: unit tests cannot be run with sh (alpine dosn't have bash by default)

acceptance-alpine-3.6: alpine-3.6
	$(call title,alpine:3.6 - acceptance)
	$(DOCKER_RUN) $(ALPINE_IMAGE) \
		sh -c '$(WGET_INSTALL_SYFT) && $(ACCEPTANCE_CMD)'

alpine-3.6:
	$(call title,alpine:3.6 - build environment)
	docker build -t $(ALPINE_IMAGE) -f $(ENVS)/Dockerfile-alpine-3.6 .


## For CI ########################################################

.PHONY: cache.fingerprint
cache.fingerprint:
	$(call title,Install test fixture fingerprint)
	@find ./environments/* -type f -exec md5sum {} + | awk '{print $1}' | sort | tee /dev/stderr | md5sum | tee cache.fingerprint && echo "$(INSTALL_TEST_CACHE_BUSTER)" >> cache.fingerprint
