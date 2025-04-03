default: run-with-docs

docs:
	$(HOME)/go/bin/swag init

run-with-docs: docs
	$(HOME)/go/bin/swag init
	templ fmt .
	templ generate
	go run main.go

build:
	$(HOME)/go/bin/swag init
	templ fmt .
	templ generate
	go build -o ./tmp/main .
