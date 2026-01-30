#!/bin/bash
# 构建多平台二进制

set -e

OUTPUT_DIR="npm-package/bin"
mkdir -p $OUTPUT_DIR

echo "Building for multiple platforms..."

# macOS
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/certctl-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o $OUTPUT_DIR/certctl-darwin-arm64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/certctl-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o $OUTPUT_DIR/certctl-linux-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o $OUTPUT_DIR/certctl-windows-amd64.exe .

echo "Build complete!"
ls -la $OUTPUT_DIR/
