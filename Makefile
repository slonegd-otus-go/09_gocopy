all:
	gofmt -w .
	go test -cover ./...
	golangci-lint run --enable-all
	go build -o gocopy

install:
	go build -o $(GOPATH)/bin/gocopy
