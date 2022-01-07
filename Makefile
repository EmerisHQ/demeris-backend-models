all: test

lint:
	golangci-lint run ./...

test:
	go test -v -race ./... -cover