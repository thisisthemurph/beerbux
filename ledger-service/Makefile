.PHONY: protos

up:
	goose up

down:
	goose down

repository:
	sqlc generate

nats:
	docker run -d --name nats-server -p 4222:4222 nats:latest
