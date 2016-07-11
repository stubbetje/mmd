
mmd_src = main.go ldd/ldd.go ldd/parser.go mmd/exporter.go mmd/definition.go

run: build/mmd
	rm -Rvf /tmp/mmd-*-tmp/
	./build/mmd -output-dir /tmp/mmd-delstef-tmp definitions/*.yaml

build/mmd: $(mmd_src)
	docker run --rm -v "$(PWD)":/go/src -w /go/src golang:1.6 go get github.com/ghodss/yaml
	docker run --rm -v "$(PWD)":/go/src -w /go/src golang:1.6 go build -v -o build/mmd
