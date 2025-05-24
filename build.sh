#!/bin/bash

cd webui || exit

npm run ci
npm run build

cd ../server

go build -tags netgo -ldflags '-s -w' -o api ./cmd/api
