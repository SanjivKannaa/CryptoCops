#!/bin/sh
CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags="-w -s -v" -a -installsuffix cgo