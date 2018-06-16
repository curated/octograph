default: test

.PHONY: test run

dev:
	docker-compose up

test:
	@CONFIG=config/config.test.json go test ./...

run:
	@go run main.go -logtostderr=true
