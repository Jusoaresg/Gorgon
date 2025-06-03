.PHONY: default docs run-with-docs build

default: 
	cd assets/front/build && npm run build 
	$(MAKE) run-with-docs

docs:
	$(HOME)/go/bin/swag init

run-with-docs: docs
	go run main.go

build: docs
	$(HOME)/go/bin/swag init
	cd assets/front/build && npm run build 
	go build -o ./tmp/main .
