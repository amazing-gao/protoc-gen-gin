#! /bin/bash

go build -o protoc-gen-gin cmd/protoc-gen-gin/*
mv protoc-gen-gin ~/go/bin/protoc-gen-gin

protoc -I/usr/local/include \
  -I . \
  -I $GOPATH/src \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I $GOPATH/src/github.com/gogo/protobuf \
  --gofast_out=logtostderr=true:. \
  example/api/api.proto


protoc -I/usr/local/include \
  -I . \
  -I $GOPATH/src \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I $GOPATH/src/github.com/gogo/protobuf \
  --gin_out=logtostderr=true:. \
  example/api/api.proto

protoc -I/usr/local/include \
  -I . \
  -I $GOPATH/src \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I $GOPATH/src/github.com/gogo/protobuf \
  --swagger_out=logtostderr=true:. \
  example/api/api.proto
