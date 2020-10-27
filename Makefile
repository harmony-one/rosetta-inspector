.PHONY: build test assets install setup docker-build

PROJECT      ?= rosetta-inspector
GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GO_VERSION   ?= $(shell go version | awk {'print $$3'})
DOCKER_IMAGE ?= figmentnetworks/${PROJECT}
DOCKER_TAG   ?= latest

build: assets
	go build

setup:
	go get -u github.com/jessevdk/go-assets-builder

assets:
	go-assets-builder static -p static -o static/assets.go

test:
	go test -race -cover ./...

install:
	go install

docker-build:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} -f Dockerfile .