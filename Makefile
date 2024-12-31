.PHONY: build run

build:
	go build -v ./cmd/desky-backend/...


run: del build
	./desky-backend -log logging

clean:
	go mod tidy
	go clean

start:
	./desky-backend -log logging

del:
	rm ./desky-backend* || echo "file didn't exists"
	rm ./trace*         || echo "file didn't exists"


build-test:
	go build -v ./cmd/desky-backend-test/...

test: build-test
	./desky-backend-test

.DEFAULT_GOAL := run