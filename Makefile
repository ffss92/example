.PHONY: build run watch

build:
	go build -o bin/api cmd/api/*

run: build
	bin/api

watch:
	reflex -r '\.go$$' -d none -s make run
