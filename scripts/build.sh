#!/bin/bash

first_param=$1

go build -o "./build/${first_param}" "./cmd/${first_param}"