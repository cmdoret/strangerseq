GOFILES = $(shell find . -name '*.go')
VERSION=""

default: release 

build/:
	mkdir -p $@


build/strangerseq-windows.exe: $(GOFILES) build/
	go get -d .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

build/strangerseq-osx: $(GOFILES) build/
	go get -d .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

build/strangerseq-linux: $(GOFILES) build/
	go get -d .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

.PHONY: release
release: build/strangerseq-linux build/strangerseq-windows.exe build/strangerseq-osx
