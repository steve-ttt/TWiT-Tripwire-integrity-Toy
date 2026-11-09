.PHONY: build test smoke

build:
	go build -o bin/tripwire .

test: build
	go test -v ./...
	./smoke_test.sh