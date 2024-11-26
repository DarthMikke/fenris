#! /bin/sh

cd /app
go install github.com/air-verse/air@latest
go mod download
go mod tidy
air -c .air.toml
