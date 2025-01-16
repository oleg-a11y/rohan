run:
	go run cmd/main.go

build:
	go build -o bid/app.exe cmd/main.go

install: 
	go mod tidy

.PHONY: run build install