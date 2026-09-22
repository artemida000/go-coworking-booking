DB_URL=postgres://postgres:postgres@localhost:5432/coworking_db?sslmode=disable

.PHONY: run build docker-up docker-down migrate-up migrate-down test

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

test:
	go test -v ./...