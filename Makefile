.PHONY: all doc test integration-test coverage

all: test integration-test coverage doc

test:
	go test -v ./...

integration-test:
	go test -tags=integration

coverage:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out
