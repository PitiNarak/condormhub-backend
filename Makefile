build:
	go build -o bin/server ./cmd/server/main.go

run: 
	go run ./cmd/server/main.go

watch:
	reflex -r '\.go$$' -s -- sh -c "make"

migrate:
	go run cmd/migrate/main.go

clean:
	rm -rf bin/server

deps:
	go mod tidy

.DEFAULT_GOAL = run