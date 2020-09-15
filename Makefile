export GOFLAGS = -mod=vendor

.PHONY: up_all
up_all:
	brew services start nginx
	@make up

.PHONY: up down
up:
	docker-compose up -d

down:
	docker-compose down

.PHONY: run
run:
	go run main.go

.PHONY: lint test tidy vendor
lint:
	golangci-lint run

test:
	go test -cover ./...

tidy:
	go mod tidy

vendor:
	go mod vendor
