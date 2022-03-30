#!/usr/bin/env sh

echo "Cleanup"
make clean

echo "Build binary"
make

echo "Run Web App"
bin/web-app-go run
