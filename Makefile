.PHONY: doc

doc:
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > doc/package.md 

integration:
	go test -tags=integration