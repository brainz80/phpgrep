.PHONY: test lint docker-test docker-lint

DOCKER_IMAGE ?= golang:1.21
DOCKER_OPTS = --rm -v "$(CURDIR):/app" -w /app
DOCKER_RUN = docker run $(DOCKER_OPTS) $(DOCKER_IMAGE)

GOPATH_DIR=`go env GOPATH`

export GO111MODULE := on

test:
	go test -v -race -count=1 -coverprofile=coverage.out ./...

GOLANGCI_LINT_VERSION := v1.62.2

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin $(GOLANGCI_LINT_VERSION)
	$(GOPATH_DIR)/bin/golangci-lint run -v

docker-test:
	$(DOCKER_RUN) make test

docker-lint:
	$(DOCKER_RUN) make lint
