.PHONY: default docs run-with-docs build

TAG ?= latest

default:
	$(MAKE) run-with-docs

docs:
	swag init

run-with-docs: docs
	go run main.go

build: docs
	mkdir -p ./tmp
	CGO_ENABLED=0 go build -o ./tmp/main .

cross-build: docs
	set -e; \
	mkdir -p ./tmp; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./tmp/main_linux_amd64 .; \
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./tmp/main_windows_amd64.exe .; \
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./tmp/main_darwin_arm64 .


test: 
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf ./tmp
	rm -rf ./configs
	rm -rf ./downloads

check: test lint

docker-build:
	docker build -t jusoares/gorgon:$(TAG) .
