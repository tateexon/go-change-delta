tidy:
	go mod tidy
lint:
	golangci-lint run

lint_fix:
	golangci-lint run --fix

test:
	go test -timeout 5m -cover -covermode=count ./...

build:
	go build -o go-change-delta main.go

eat_dogfood:
	go test -timeout 5m -cover -covermode=count $(shell go run main.go -b=origin/main -l=0)

typos:
	typos
