#!/bin/sh

go build -o ./bin/isoserver cmd/isoserver/main.go
go build -o ./bin/isoctl cmd/isoctl/main.go