# 'run' target: Runs the project directly in development mode
run:
	go run cmd/main.go

# 'build' target: Builds the binary and outputs it as app.exe in the 'bin' directory
build:
	go build -o bin/app.exe cmd/main.go 

# 'install' target: Installs the dependencies by running 'go mod tidy'
install:
	go mod tidy 

# Declare these targets as 'phony' to avoid conflicts with files or directories of the same name
.PHONY: run build install