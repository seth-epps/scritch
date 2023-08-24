generate:
	go generate ./...

build: generate
	go build -o dist/scritch
