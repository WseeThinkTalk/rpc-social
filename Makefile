## rpc 生成命令
rpc:
	goctl rpc protoc social.proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=./social --go-grpc_out=./social --zrpc_out=./ -m --verbose --style=go_zero

run:
	go run social.go

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/rpc-social social.go
