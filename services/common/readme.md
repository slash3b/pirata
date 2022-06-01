
#### how to add `local` package 
- `go mod edit -require=common@v0.0.0  -replace=common@v0.0.0=../common`
- check tutorial https://go.dev/doc/tutorial/call-module-code

#### generate proto files
protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative imdb.proto
#### generate grpc files



#### official documentation
https://grpc.io/docs/languages/go/quickstart/