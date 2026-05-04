.PHONY: build run test up down swagger

build:
	go build -o bin/subscription-service ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

swagger:
	swag init -g cmd/api/main.go -o docs

up:
	docker-compose up --build

down:
	docker-compose down

clean:
	rm -rf bin/