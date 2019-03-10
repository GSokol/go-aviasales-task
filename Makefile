#
# Makefile
# Grigorii Sokolik, 2019-02-25 17:41
#

GO?=go
TEST_FLAGS?=
TEST_FLAGS+= -race -v -vet all -cover
DEP?=${GOPATH}/bin/dep
SOURCE_FILES=$(shell find ./cmd ./pkg ./internal -name "*.go")
DEP_FILE=Gopkg.toml
DEP_LOCK_FILE=Gopkg.lock
BINARY_NAME=./bin/go-aviasales-task
GO_BUILD_ENV_VARS?=CGO_ENABLED=0 GOOS=linux
GO_BUILD_ENV_FLAGS?=-a -installsuffix cgo
GO_BUILD_ENV_FLAGS+= -o ${BINARY_NAME}
ENTRYPOINT_PACKAGE=./cmd/service

${BINARY_NAME}: ${DEP_LOCK_FILE}
	${GO_BUILD_ENV_VARS} ${GO} build ${GO_BUILD_ENV_FLAGS} ${ENTRYPOINT_PACKAGE}

.PHONY: test

test: ${DEP_LOCK_FILE}
	${GO} test ${TEST_FLAGS} ./cmd/... ./internal/... ./pkg/...



${DEP_LOCK_FILE}: ${SOURCE_FILES} ${DEP_FILE}
	${DEP} ensure

# vim:ft=make
#
