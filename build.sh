#!/usr/bin/bash

cd frontend || exit
bun install
bun run build

cd ..
go mod vendor
go build -o bin/mcquery main.go
