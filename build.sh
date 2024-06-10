#!/bin/bash

# Create a directory for the build output
mkdir -p dist

# Build the Go binary
go build -o ./dist/processamento-imagens ./cmd/

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ./dist/processamento-imagens.exe ./cmd/

# Copy static files and other directories
cp -r ./static ./dist/
cp -r ./storage ./dist/
cp -r ./view ./dist/
