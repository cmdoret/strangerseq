GOFILES = $(shell find . -name '*.go')
VERSION=""

default: release 

build/:
	mkdir -p $@

.PHONY: deps
deps:
	go get -d .

build/strangerseq-windows.exe: $(GOFILES) build/ deps
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

build/strangerseq-osx: $(GOFILES) build/ deps
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

build/strangerseq-linux: $(GOFILES) build/ deps
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $@ .

.PHONY: release
release: build/strangerseq-linux build/strangerseq-windows.exe build/strangerseq-osx

.PHONY test
test: deps
	go test ./kmers
