#!/bin/bash

echo "Running: GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o server.exe server.go"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o server.exe server.go

echo "Running: go build -o server-linux server.go"
go build -o server-linux server.go
