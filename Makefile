include .env
export

ref:
	@echo "Tidy and vendor"
	@go mod tidy && go mod vendor
PHONY: ref

run:
	@echo "Running the application"
	@echo "running on: $$ENV"
	@echo "DB: $$DB_HOST"
	@echo "QUEUE: $$QUEUE_HOST"
	@go run .
PHONY: run