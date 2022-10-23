include .env

.PHONY: help

help: ## Output help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run billing app along with PostgreSQL server using docker-compose
	docker-compose up --build -d && docker-compose logs -f

compose-down: ### Stop billing app and PostgreSQL server launched using docker-compose
	docker-compose down --remove-orphans

compose-test: ### Run integration testing in docker environment (Warning: recreates containers, do not use with persistent db storage)
	docker-compose rm -f
	docker-compose --profile testing up --build  --abort-on-container-exit --exit-code-from newman

run: build ### Run billing app locally
	cd cmd && ./billing

build: ### Build billing app locally
	go mod download && go build -o ./cmd/billing ./cmd/main.go
