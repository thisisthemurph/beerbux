.PHONY: protos

protos:
	protoc -I protos/ protos/transaction.proto --go_out=.
	protoc -I protos/ protos/transaction.proto --go-grpc_out=.
