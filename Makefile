.PHONY: build run

# ========= Vars definitions =========

app = application

# ========= Prepare commands =========

tidy:
	go mod tidy
	go clean

del:
	rm ./$(app)* || echo "file didn't exists"
	rm ./trace*  || echo "file didn't exists"

# ========= Compile commands =========

build:
	go build -v ./cmd/$(app)/...

run: del build
	./$(app)

cross-compile:

	GOOS=linux
	go build -o ./dist/linux/$(app) -v ./cmd/$(app)/...

	GOOS=windows
	go build -o ./dist/windows/$(app).exe -v ./cmd/$(app)/...

# ========= Documentation generate =========

swag:
	swag init -g ./cmd/$(app)/main.go

doc:
	godoc -http localhost:3000

.DEFAULT_GOAL := run