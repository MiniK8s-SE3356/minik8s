CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/scheduler ./pkg/scheduler/main
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/apiserver ./cmd/apiserver
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubectl ./cmd/kubectl
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/controller ./cmd/controller
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./build/kubelet ./pkg/kubelet/main/