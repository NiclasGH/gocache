dependencies:
	go mod tidy

run:
	go run cmd/gocache/main.go

build:
	go build cmd/gocache/main.go

install:
	go install cmd/gocache/main.go

fmt:
	go fmt ./...

check:
	go build -o /dev/null cmd/gocache/main.go
