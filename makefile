default: run-with-docs

docs:
	$(HOME)/go/bin/swag init

run-with-docs: docs
	$(HOME)/go/bin/swag init
	go run main.go

