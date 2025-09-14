.PHONY: install protos

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build:
	go mod tidy

run:
	cd server && go run main.go

protos:
	protoc --proto_path=protos --go_out=. --go-grpc_out=. $(shell ls protos/*.proto)
