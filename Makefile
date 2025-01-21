run:
	go run cmd/main.go

build:
	go build -o bin/app.exe cmd/main.go

install: 
	go mod tidy

.PHONY: run build install