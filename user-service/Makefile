.PHONY: protos

protos:
	protoc -I protos/ protos/user.proto --go_out=.
	protoc -I protos/ protos/user.proto --go-grpc_out=.
