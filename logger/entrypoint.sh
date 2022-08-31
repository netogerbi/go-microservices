#!/bin/sh

apt-get update
apt-get install -y protobuf-compiler

export GO111MODULE=on
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=paths=source_relative:. ./logs/logs.proto --go-grpc_out=. ./logs/logs.proto

# the command configured in CMD on dockerfile will run below in "$@"
exec "$@"
