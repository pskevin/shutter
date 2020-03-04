go run main.go CreateNode 1 1000
go run main.go CreateNode 2 500
go run main.go Send 1 2 300
go run main.go Send 2 1 100
go run main.go Send 1 2 400
go run main.go Receive 2 1
go run main.go BeginSnapshot 1
go run main.go Receive 2 1
go run main.go Send 1 2 100
go run main.go Receive 2 1
go run main.go ReceiveAll
go run main.go CollectState
go run main.go PrintSnapshot
go run main.go KillAll