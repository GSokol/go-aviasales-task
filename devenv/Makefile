#
# Makefile
# Grigorii Sokolik, 2019-02-26 09:53
#

# Constants
DOCKER_IMAGE_NAME=go-aviasales-task

# Functions
slashed = $(shell echo "$(1)" | sed 's/^\s\+//' | sed 's/\s\+$$//' | sed 's/\([^\/]\)$$/\1\//')
sed_escape = $(subst /,\/,$(subst .,\.,$(1)))

# Variablse
DOCKER?=docker
DOCKER_FILE?=Dockerfile
DOCKER_REGISTRY?=
DOCKER_PROJECT?=gsokol
VALUES_FILE?=values.yaml

# Templates
GIT_REV:=$(shell git rev-parse --short HEAD)
GIT_TAGS:=$(shell git tag --list --points-at=HEAD | xargs echo)
DATE:=$(shell date +%s)
KUBER_NAMESPACE?=${DOCKER_PROJECT}
DOCKER_IMAGE_TAG?=${DATE}.${GIT_REV}
TEST_IMAGE_TAG?=test-${DOCKER_IMAGE_TAG}
DOCKER_IMAGE_FULLNAME:=$(call slashed,${DOCKER_REGISTRY})$(call slashed,${DOCKER_PROJECT})${DOCKER_IMAGE_NAME}
DOCKER_BUILD_TAGS:=$(shell echo ' ${DOCKER_IMAGE_TAG} ${GIT_TAGS}' | sed 's/ \(\S\+\)/ $(call sed_escape,${DOCKER_IMAGE_FULLNAME}):\1/g')
DOCKER_BUILD_TAGS_CMD:=$(shell echo '${DOCKER_BUILD_TAGS}' | sed 's/ \(\S\+\)/ -t \1/g')
HELM_RELEASE_NAME?=${DOCKER_IMAGE_NAME}
HELM_VERSION?=${DATE}
PWD:=$(shell pwd)

.PHONY: test
test: build-test-image
	docker run --rm -it \
		-v ${PWD}:/go/src/github.com/GSokol/go-aviasales-task \
		${DOCKER_IMAGE_NAME}-test

.PHONY: build-test-image
build-test-image:
	docker build \
		-f ${PWD}/devenv/Dockerfile \
		-t ${DOCKER_IMAGE_NAME}-test .

.PHONY: build
build:
	docker build \
		-f ${DOCKER_FILE} \
		${DOCKER_BUILD_TAGS_CMD} \
		.

.PHONY: push
push: build
	docker push \
		${DOCKER_BUILD_TAGS}

helm/requirements.lock: helm/requirements.yaml
	helm dependency update ./helm

.PHONY: deploy-k8s
ifeq (${SHOULD_PUSH},1)
deploy-k8s: push helm/requirements.lock
else
deploy-k8s: build helm/requirements.lock
endif
	helm upgrade \
		${HELM_RELEASE_NAME} \
		./helm \
		--install \
		--namespace ${KUBER_NAMESPACE} \
		--version ${HELM_VERSION} \
		--set helmChart.image.registry=$(call slashed,${DOCKER_REGISTRY}) \
		--set helmChart.image.project=$(call slashed,${DOCKER_PROJECT}) \
		${HELM_SETTINGS}--set helmChart.image.tag=${DOCKER_IMAGE_TAG} \
		-f helm/${VALUES_FILE}

.PHONY: build-debug-image
build-debug-image:
	docker build \
		-f ${PWD}/devenv/Dockerfile.debug \
		-t ${DOCKER_IMAGE_NAME}:debug .

.PHONY: deploy-debug-image
deploy-debug-image: build-debug-image
	helm upgrade \
		${HELM_RELEASE_NAME} \
		./helm \
		--install \
		--namespace ${KUBER_NAMESPACE} \
		--version ${HELM_VERSION} \
		--set helmChart.image.registry="" \
		--set helmChart.image.project="" \
		${HELM_SETTINGS}--set helmChart.image.tag=debug \
		--set helmChart.debug=true \
		--set helmChart.deployment.replicas=1 \
		-f helm/${VALUES_FILE}

# vim:ft=make
#
