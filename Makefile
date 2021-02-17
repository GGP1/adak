test:
	go test ./...

test-race:
	go test ./... -race

run:
	go run cmd/main.go

run-api:
	go run api/swagger/main.go

docker-build-api:
	docker build -t adak-api api/swagger

docker-run-api:
	docker run --rm -p 8080:8080 adak-api