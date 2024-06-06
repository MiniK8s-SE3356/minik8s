#!/bin/bash

PATH=$PATH:/usr/local/go/bin
go mod tidy
# 安装mock工具
go get -u github.com/golang/mock/mockgen
go test ./...