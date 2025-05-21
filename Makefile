build:
	@go fmt ./...
	@go build -o bin/pulse

run: build
	@./bin/pulse

test:
	@go test -v ./...

up:
	@goose -dir ./db/migrations/ -no-versioning up

seed:
	@goose -dir ./db/seed/ -no-versioning up

down:
	@goose -dir ./db/migrations/ -no-versioning down
