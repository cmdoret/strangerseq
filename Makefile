GOFILES = $(shell find . -name '*.go')
VERSION=""

default: build

workdir:
	mkdir -p build/


.PHONY: windows
build/strangerseq: $(GOFILES)
	go get -d .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o build/strangerseq-windows.exe .

.PHONY: osx
build/strangerseq: $(GOFILES)
	go get -d .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o build/strangerseq-osx .

.PHONY: linux
build/strangerseq: $(GOFILES)
	go get -d .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o build/strangerseq-linux .

.PHONY: release
build: windows linux osx
