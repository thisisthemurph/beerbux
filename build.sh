#!/bin/bash

cd webui || exit

npm ci
npm run build

cd ../server

go build -tags netgo -ldflags '-s -w' -o api ./cmd/api
