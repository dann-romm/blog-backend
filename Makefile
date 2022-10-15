include .env
export

# rewrite PH_HOST env variable



.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

compose-re: compose-down compose-up ### Restart docker-compose
.PHONY: compose-re

docker-rm-volume: ### remove docker volume
	docker volume rm pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' -c --ignored | xargs hadolint
.PHONY: linter-hadolint

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'blog_backend'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	echo "y" | migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' down
.PHONY: migrate-down
