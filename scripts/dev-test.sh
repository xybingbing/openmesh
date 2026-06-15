#!/bin/sh
set -eu

go test ./...
go build -o ./bin/openmesh ./cmd/openmesh
