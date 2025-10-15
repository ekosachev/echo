.PHONY: run test docker-up docker-down

run:
	ENV=development HTTP_PORT=8080 JWT_SECRET=supersecretjwtkey DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=appdb DB_SSLMODE=disable TIME_ZONE=UTC go run ./cmd/server

test:
	go test ./...

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down -v
