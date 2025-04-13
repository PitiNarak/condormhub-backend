build:
	go build -o bin/server ./main.go

run:
	go run ./main.go

dev:
	air

migrate:
	go run cmd/migrate/main.go

clean:
	rm -rf bin/server

deps:
	go mod tidy

gen-docs:
	swag init -v3.1 -o docs -g main.go

lint:
	golangci-lint run

.DEFAULT_GOAL = run
