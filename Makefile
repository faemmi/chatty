.PHONY: install protos

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

protos:
	protoc --proto_path=protos --go_out=. --go-grpc_out=. $(shell ls protos/*.proto)

build: protos
	go mod tidy
	go -C server build

run: build
	go run server/main.go

test: build
	go test -coverprofile=coverage.out -v ./server/...
	go tool cover -func=coverage.out
