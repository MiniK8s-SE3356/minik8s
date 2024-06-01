#!/bin/bash

cd /home/xubbbb/Code/CloudOS/minik8s/
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./pkg/gpu/jobserver/imagebuilding/jobserver ./pkg/gpu/jobserver/main/
docker rmi -f jobserver-image levixubbbb/jobserver-image
cd /home/xubbbb/Code/CloudOS/minik8s/pkg/gpu/jobserver/imagebuilding
docker build -t jobserver-image .
docker tag jobserver-image:latest levixubbbb/jobserver-image:latest
docker push levixubbbb/jobserver-image:latest