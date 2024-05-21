go build -o scheduler ./pkg/scheduler/main
go build -o apiserver ./pkg/apiserver
go build -o kubectl ./pkg/kubectl
go build -o controller ./cmd/controller