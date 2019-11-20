.PHONY: all doc test integration-test coverage

all: test integration-test coverage doc

doc:
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > doc/package.md 

test:
	go test -v ./...

integration-test:
	go test -tags=integration

coverage:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out