CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/scheduler_static ./pkg/scheduler/main
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/apiserver_static ./cmd/apiserver
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubectl_static ./cmd/kubectl
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/controller_static ./cmd/controller
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubelet_static ./pkg/kubelet/main/