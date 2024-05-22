go build -o build/scheduler ./pkg/scheduler/main
go build -o build/apiserver ./cmd/apiserver
go build -o build/kubectl ./cmd/kubectl
go build -o build/controller ./cmd/controller