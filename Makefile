default: test
.PHONY: dist test run

dist:
	@docker build -t octograph .

test:
	@go test ./...

run:
	@go run main.go
