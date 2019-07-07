all:
	gofmt -w .
	go test -cover ./...
	golangci-lint run --enable-all
	go build

