default: test

.PHONY: dev test index delete

dev:
	docker-compose up

test:
	@CONFIG=config/config.test.json go test ./...

index:
	@go run main.go -logtostderr=true

delete:
	@go run main.go -logtostderr=true -process=delete
