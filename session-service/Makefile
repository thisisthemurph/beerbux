.PHONY: protos

up:
	goose up

down:
	goose down

repository:
	sqlc generate

protos:
	protoc -I protos/ protos/session.proto --go_out=.
	protoc -I protos/ protos/session.proto --go-grpc_out=.
	protoc -I protos/ protos/history.proto --go_out=.
	protoc -I protos/ protos/history.proto --go-grpc_out=.

nats:
	docker run -d --name nats-server -p 4222:4222 nats:latest
