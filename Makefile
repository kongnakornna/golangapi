.PHONY: clean tidy download verify migrate swag vendor test run

clean:
	go clean -cache
	go clean -modcache

tidy:
	go mod tidy

download:
	go mod download

verify:
	go mod verify

migrate:
	go run cmd/api/main.go migrate

swag:
	swag init -g cmd/api/main.go

vendor:
	go mod vendor

test:
	go test ./...

run: clean tidy download verify migrate swag vendor test
	air