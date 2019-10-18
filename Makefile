GOFILES = $(shell find . -name '*.go')

default: build

workdir:
	mkdir -p workdir

build: workdir/strangerseq

workdir/strangerseq: $(GOFILES)
	go get -d .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o workdir/strangerseq .
