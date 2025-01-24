.PHONY: build run

app = desky-backend
test-app = desky-backend-test

build:
	go build -v ./cmd/$(app)/...


run: del build
	./$(app) -log logging

clean:
	go mod tidy
	go clean

start:
	./$(app) -log logging

del:
	rm ./$(app)* || echo "file didn't exists"
	rm ./trace*         || echo "file didn't exists"


build-test:
	go build -v ./cmd/$(app-test)/...

test: build-test
	./$(app-test)

swag:
	swag init -g ./cmd/$(app)/main.go

doc:
	godoc -http localhost:3000

.DEFAULT_GOAL := run