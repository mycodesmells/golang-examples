#!/bin/bash
echo "Building golang binary"
GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o app ../../web
