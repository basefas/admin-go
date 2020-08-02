.PHONY: run build

run:
	go run ./cmd/app/main.go
build:
	go build -o ./build/bin/admin-go -v ./cmd/app/main.go