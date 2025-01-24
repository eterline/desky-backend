.PHONY: build run

app = application

build:
	go build -v ./cmd/$(app)/...


run: del build
	./$(app)

clean:
	go mod tidy
	go clean

start:
	./$(app)

del:
	rm ./$(app)* || echo "file didn't exists"
	rm ./trace*         || echo "file didn't exists"

swag:
	swag init -g ./cmd/$(app)/main.go

doc:
	godoc -http localhost:3000

.DEFAULT_GOAL := run