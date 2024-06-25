#!/usr/bin/bash

cd frontend || exit
bun install
bun run build

cd ..
go mod vendor
go build -o build/mcquery main.go
