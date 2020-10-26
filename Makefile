.PHONY: build test assets install

build: assets
	go build

assets:
	go-assets-builder static -p static -o static/assets.go

test:
	go test -race -cover ./...

install:
	go install
