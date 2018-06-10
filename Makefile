default: test

.PHONY: test run

test:
	@go test ./...

run:
	@go run main.go
