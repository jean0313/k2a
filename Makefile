VERSION := $(shell git describe --always --dirty)

.PHONY: clean

all: linux windows macos

windows:
	GOOS=windows GOARCH=amd64 go build -a -ldflags="-X main.Version=$(VERSION)" -o build/windows/cli.exe

macos:
	GOOS=darwin GOARCH=amd64 go build -a -ldflags="-X main.Version=$(VERSION)" -o build/macos/cli

linux:
	GOOS=linux GOARCH=amd64 go build -a -ldflags="-X main.Version=$(VERSION)" -o build/linux/cli

clean:
	go clean
	rm -rf build/

