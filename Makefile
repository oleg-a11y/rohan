.PHONY: install
run:
	go mod tidy

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	go build -o bin/app.exe cmd/main.go