.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test ./... -race

run:
	go run cmd/main.go

.PHONY: rebuild-server
rebuild-server:
	docker compose rm -sf server && docker compose up --build --no-deps server