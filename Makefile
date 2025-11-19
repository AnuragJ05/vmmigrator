    PROTO_DIR=proto
    PROTO_FILE=$(PROTO_DIR)/vmmigrator.proto

    .PHONY: generate build run

    generate:
	protoc --go_out=paths=source_relative:proto --go-grpc_out=paths=source_relative:proto proto/vmmigrator.proto

    build: generate
	go build ./...

    run: build
	go run ./services/api-gateway/cmd
