generate:
	protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative imdb.proto
	protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative imdb.proto