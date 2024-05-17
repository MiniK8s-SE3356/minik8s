#!/bin/bash

first_param=$1
export CGO_ENABLED=0
go build -o "./build/${first_param}_static" "./cmd/${first_param}"
export CGO_ENABLED=1