.PHONY: build
build:
		go build -v ./cmd/app/

.PHONY: test
test:
		go test -v -race -timeout 30s ./...
.DEFAULT_GOAL := build

migrate:
		migrate -path ./schema -database 'postgres://postgres:straykz@localhost:5436/postgres?sslmode=disable' up
