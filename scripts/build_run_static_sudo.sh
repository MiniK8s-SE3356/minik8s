export CGO_ENABLED=0
go build -o ./build/kubeProxy ./cmd/kubeProxy/
sudo ./build/kubeProxy
export CGO_ENABLED=1