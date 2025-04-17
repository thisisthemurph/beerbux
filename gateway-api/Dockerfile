FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o gateway-api

CMD ["./gateway-api"]
