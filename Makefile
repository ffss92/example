.PHONY: build run watch

build:
	go build -o bin/api cmd/api/*

run: build
	bin/api

watch:
	reflex -r '\.go$$' -d none -s make run

test/cover:
	@mkdir -p coverage
	go test ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
