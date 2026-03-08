#!/bin/sh
docker build -t subito-assignment-test --target builder .
docker run --rm subito-assignment-test go test -v ./...

