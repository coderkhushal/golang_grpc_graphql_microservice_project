# Packages Needed
- Protobuff
- protoc-gen-go-grpc package
```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
- protoc-go-gen package
```shell
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
# Command To Generate grpc and pb files are
```shell
protoc --go_out=. --go-grpc_out=. account/account.proto
```