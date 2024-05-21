go build -o build/scheduler ./pkg/scheduler/main
go build -o build/apiserver ./pkg/apiserver
go build -o build/kubectl ./pkg/kubectl
go build -o build/controller ./cmd/controller