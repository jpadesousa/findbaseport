BINARY := findbaseport 

.PHONY: all fmt vet lint build clean

all: fmt vet lint build

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

build:
	mkdir -p bin
	go build -o bin/$(BINARY)

clean:
	rm -rf bin/
