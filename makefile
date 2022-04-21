
.PHONY: example
example:
	go build -o protoc-gen-gin cmd/protoc-gen-gin/*
	mv protoc-gen-gin $$GOROOT/bin/protoc-gen-gin
	protoc -I . -I third_party --go_out . --go_opt paths=source_relative example/api/api.proto
	protoc -I . -I third_party --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative example/api/api.proto
	protoc -I . -I third_party --gin_out . example/api/api.proto

	cd example && go run .