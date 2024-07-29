tidy:
	go mod tidy
lint:
	golangci-lint run

lint_fix:
	golangci-lint run --fix

test:
	go test -timeout 5m ./...

test_cover:
	go test -timeout 5m -cover -covermode=count ./...

build:
	go build -o go-change-delta main.go

eat_dogfood:
	go test -v -timeout 5m $(shell go run main.go -b=main)