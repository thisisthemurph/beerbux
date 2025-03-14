.PHONY: protos

protos:
	protoc -I protos/ protos/auth.proto --go_out=.
	protoc -I protos/ protos/auth.proto --go-grpc_out=.
