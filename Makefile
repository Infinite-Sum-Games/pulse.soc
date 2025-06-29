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

docker:
	@docker compose up -d

dev:
	@podman compose down
	@podman compose up -d

# Populate LIVE Updates stream with seed data
live:
	@bash scripts/stream.sh

# Populate Language Sorted-Sets with seed data
season:
	@bash scripts/sorted-set.sh

# Populate Leaderboard Sorted-Set with seed data
board:
	@bash scripts/leaderboard.sh

# Add the seed data and remove test-only seeds
prod:
	@goose up
	@goose -dir ./db/seed/ up
	@goose -dir ./db/seed/ down
