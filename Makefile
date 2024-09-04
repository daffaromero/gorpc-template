run-service:
	go run .

gen-api:
	@protoc \
    --proto_path=protobuf "protobuf/api.proto" \
    --go_out=protobuf/api --go_opt=paths=source_relative \
    --go-grpc_out=protobuf/api --go-grpc_opt=paths=source_relative