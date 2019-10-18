GOFILES = $(shell find . -name '*.go')
VERSION=""

default: build

workdir:
	mkdir -p workdir

build: workdir/strangerseq

workdir/strangerseq: $(GOFILES)
	go get -d .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o workdir/strangerseq .
