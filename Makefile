.PHONY: help run test test-v migrate fmt

GO ?= go

help:
	@echo "Available targets:"
	@echo "  run       Run the API server (go run ./cmd)"
	@echo "  test      Run tests (go test ./...)"
	@echo "  test-v    Run tests verbose (go test -v ./...)"
	@echo "  migrate   Run migrations (go run ./cmd/migrate)"
	@echo "  migrate-up Apply migrations (go run ./cmd/migrate up)"
	@echo "  migrate-down Rollback migrations (go run ./cmd/migrate down)"
	@echo "  fmt       Format Go files (gofmt -w .)"

run:
	$(GO) run ./cmd

test:
	$(GO) test ./...

test-v:
	$(GO) test -v ./...

migrate:
	$(GO) run ./cmd/migrate


start-mysql:
	docker start myproject-mysql || docker run -d --name myproject-mysql \
		-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
		-e MYSQL_DATABASE=myapp \
		-p 3306:3306 mysql:8.0
migrate-up:
	$(GO) run ./cmd/migrate up

migrate-down:
	$(GO) run ./cmd/migrate down

fmt:
	gofmt -w .
