.PHONY: build build-plugin test

build:
	go build -o loglint ./main.go
	
toLint:
	golangci-lint run

build-plugin:
	go build -buildmode=plugin -o loglint.so ./plugin

test:
	go test -v ./...

vet:
	go vet -vettool=./loglint ./...
