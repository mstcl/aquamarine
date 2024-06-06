.DEFAULT_GOAL := build
.PHONY: lint tidy build

golangci-lint:
	golangci-lint run

tidy:
	go mod tidy

build:
	go build -gcflags "-l" -ldflags "-w -s" -o aquamarine main.go

install:
	go install -gcflags "-l" -ldflags "-w -s" .
