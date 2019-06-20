#! /bin/bash

go build -o protoc-gen-gin cmd/protoc-gen-gin/*
mv protoc-gen-gin ~/go/bin/protoc-gen-gin

protoc -I/usr/local/include \
  -I . \
  -I third_party \
  --gofast_out=logtostderr=true:. \
  example/api/api.proto


protoc -I/usr/local/include \
  -I . \
  -I third_party \
  --gin_out=logtostderr=true:. \
  example/api/api.proto

protoc -I/usr/local/include \
  -I . \
  -I third_party \
  --swagger_out=logtostderr=true:. \
  example/api/api.proto
