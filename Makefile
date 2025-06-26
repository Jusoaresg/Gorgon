.PHONY: default docs run-with-docs build

default: 
	cd assets/front/build && npm run build 
	$(MAKE) run-with-docs

docs:
	$(HOME)/go/bin/swag init

run-with-docs: docs
	go run main.go

build: docs
	mkdir -p ./tmp
	$(HOME)/go/bin/swag init
	cd assets/front/build && npm run build 
	go build -o ./tmp/main .

cross-build:
	mkdir -p ./tmp
	$(HOME)/go/bin/swag init
	cd assets/front/build && npm run build 
	GOOS=linux GOARCH=amd64 go build -o ./tmp/main_linux_amd64 .
	GOOS=windows GOARCH=amd64 go build -o ./tmp/main_windows_amd64.exe .
	GOOS=darwin GOARCH=arm64 go build -o ./tmp/main_darwin_arm64 .


test: 
	go test ./...
