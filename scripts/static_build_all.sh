CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/scheduler ./cmd/scheduler
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/apiserver ./cmd/apiserver
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubectl ./cmd/kubectl
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/controller ./cmd/controller
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubelet ./cmd/kubelet
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/serverless ./cmd/serverless
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubeProxy ./cmd/kubeProxy
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/jobmanager ./cmd/jobmanager