#!/bin/bash

PATH=$PATH:/usr/local/go/bin
go mod tidy
# go test ./...