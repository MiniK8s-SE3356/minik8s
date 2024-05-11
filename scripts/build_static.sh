export CGO_ENABLED=0
go build -o ./build/kubeProxy ./cmd/kubeProxy/
export CGO_ENABLED=1