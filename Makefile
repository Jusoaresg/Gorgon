default: run-with-docs

docs:
	$(HOME)/go/bin/swag init

templ: 
	@templ generate --watch --proxy=http://localhost:8080

run-with-docs: docs
	$(HOME)/go/bin/swag init
	templ fmt .
	templ generate
	go run main.go

build:
	$(HOME)/go/bin/swag init
	templ generate
	go build -o ./tmp/main .
