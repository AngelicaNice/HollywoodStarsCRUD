build:
	docker-compose build app

run:
	docker-compose up app

migrate:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go

lint:
	golangci-lint run