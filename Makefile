build:
	@go fmt ./...
	@go build -o bin/pulse

run: build
	@./bin/pulse

test:
	@go test -v ./...

