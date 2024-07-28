tidy:
	go mod tidy
lint:
	golangci-lint run

test_unit: go_mod
	go test -timeout 5m -json -cover -covermode=count -coverprofile=unit-test-coverage.out 2>&1 | tee /tmp/gotest.log | gotestloghelper -ci