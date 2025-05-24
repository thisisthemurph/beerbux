#!/bin/bash

cd webui || exit

export VITE_API_BASE_URL=/api
export VITE_SSE_BASE_URL=/api/events
export VITE_TEST="this is only a test"

echo "API_BASE_URL: $VITE_API_BASE_URL"

npm ci
npm run build

cd ../server

go build -tags netgo -ldflags '-s -w' -o api ./cmd/api
