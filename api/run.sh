#! /bin/sh

cd /app
go install github.com/air-verse/air@latest
go mod download
air -c .air.toml
