.PHONY: all doc test integration-test coverage

all: test integration-test coverage doc

doc:
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > doc/package.md 

test:
	gotest -v ./...

integration-test:
	gotest -tags=integration

coverage:
	gotest -coverprofile=coverage.out 
	go tool cover -html=coverage.out