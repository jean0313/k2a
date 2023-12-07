VERSION = $(shell git describe --always --dirty)

linux: main.go
    GOOS=linux \
      go build -a -ldflags="-X main.Version=$(VERSION)" -o build/cli

windows: main.go
	GOOS=windows GOARCH=amd64\
      go build -a -ldflags="-X main.Version=$(VERSION)" -o build/cli.exe

macos: main.go
	GOOS=darwin GOARCH=amd64\
      go build -a -ldflags="-X main.Version=$(VERSION)" -o build/cli