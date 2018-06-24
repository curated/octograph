default: test

.PHONY: test index delete push release

test:
	@CONFIG=config/test.config.json GOCACHE=off go test ./...

index:
	@go run main.go -logtostderr=true

delete:
	@go run main.go -logtostderr=true -process=delete

push:
	@heroku container:push worker --app=octograph

release:
	@heroku container:release worker --app=octograph
