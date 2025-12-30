# Proje Değişkenleri
APP_NAME=ecommerce-api
DOCKER_COMPOSE=docker-compose.yml

.PHONY: up down restart build test migrate-up migrate-down

up:
	docker-compose up -d --build

down:
	docker-compose down

restart:
	docker-compose restart app

logs:
	docker-compose logs -f app

test:
	go test -v ./...

# Örnek: Migrate işlemleri için (golang-migrate yüklü olmalı)
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/ecommerce?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/ecommerce?sslmode=disable" down