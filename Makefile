.PHONY: build test assets install setup

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
