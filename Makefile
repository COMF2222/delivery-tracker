include .env

DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o docs/ --parseDependency

.PHONY: swagger-fmt
swagger-fmt:
	swag fmt

# Пересоздавать документацию перед каждым запуском
.PHONY: run
run: swagger
	go run ./cmd/api/main.go