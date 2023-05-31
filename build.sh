#!/bin/sh

go generate

env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/tunnel
env GOOS=darwin GOARCH=amd64 go build -o build/macos-amd64/tunnel
env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/tunnel
