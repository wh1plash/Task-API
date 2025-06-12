build:
	@go build -o bin/app main.go

run: build
	@./bin/app

docker: build
	@docker	compose up --build -d

seed:
	@go run scripts/seed.go

.PHONY: build	