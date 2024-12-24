.PHONY: build run

build:
	go build -v ./cmd/desky-backend/...


run: del build
	./desky-backend

clean:
	go mod tidy
	go clean

start:
	./desky-backend -log ./logging

del:
	rm ./desky-backend || echo "file didn't exists"

.DEFAULT_GOAL := run