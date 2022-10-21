.PHONY: build
build:
		go build -v ./cmd/app/

.PHONY: test
test:
		go test -v -race -timeout 30s ./...
.DEFAULT_GOAL := build

.PHONY: migrate
migrate:
		migrate -path ./schema -database 'postgres://postgres:straykz@localhost:5436/postgres?sslmode=disable' up
.PHONY: createDocker
createDocker:
		docker run --name=avitoDB -e POSTGRES_PASSWORD=straykz -p 5436:5432 -d --rm postgres
