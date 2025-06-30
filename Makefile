.PHONY: default docs run-with-docs build

default: 
	cd assets/front/build && npm run build 
	$(MAKE) run-with-docs

docs:
	swag init

run-with-docs: docs
	go run main.go

build-front:
	cd assets/front && npm run build 

build: docs build-front
	mkdir -p ./tmp
	go build -o ./tmp/main .

cross-build: docs build-front
	set -e; \
	mkdir -p ./tmp; \
	GOOS=linux GOARCH=amd64 go build -o ./tmp/main_linux_amd64 .; \
	GOOS=windows GOARCH=amd64 go build -o ./tmp/main_windows_amd64.exe .; \
	GOOS=darwin GOARCH=arm64 go build -o ./tmp/main_darwin_arm64 .


test: 
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf ./tmp

check: build-front test lint
