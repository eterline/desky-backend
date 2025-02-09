.PHONY: build run

# ========= Vars definitions =========

app = application
db = migrator


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
	go build -v ./cmd/$(db)/...

run: del build
	./$(app)

cross-compile:

	GOOS=linux
	go build -o ./dist/linux/$(app) -v ./cmd/$(app)/...
	go build -o ./dist/linux/$(db) -v ./cmd/$(db)/...

	GOOS=windows
	go build -o ./dist/windows/$(app).exe -v ./cmd/$(app)/...
	go build -o ./dist/windows/$(db).exe -v ./cmd/$(db)/...

# ========= Documentation generate =========

swag:
	swag init -g ./cmd/$(app)/main.go

doc:
	godoc -http localhost:3000

.DEFAULT_GOAL := run