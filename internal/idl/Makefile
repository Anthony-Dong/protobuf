.PHONY: all build go

all: build
GOPATH := $(shell go env GOPATH)

## install  protoc & protoc-gen-go & protoc-gen-go-grpc
## wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.0/protoc-3.19.0-osx-x86_64.zip
## go get -v google.golang.org/protobuf/cmd/protoc-gen-go
## go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc

build:
	protoc --proto_path desc \
	--proto_path . \
	--plugin=protoc-gen-go=${HOME}/go/bin/protoc-gen-go \
	--go_out=${HOME}/go/src \
	--plugin=protoc-gen-go-grpc=${HOME}/go/bin/protoc-gen-go-grpc \
	--go-grpc_out=${HOME}/go/src \
	`find . -name '*.proto'`