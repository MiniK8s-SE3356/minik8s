#!/bin/bash

docker rmi -f levixubbbb/jobserver-image
etcdctl del "/minik8s/gpujob" --prefix
etcdctl del "/minik8s/pod" --prefix
etcdctl del "/minik8s/node" --prefix