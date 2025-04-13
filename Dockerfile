FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4

COPY . .

RUN	swag init -v3.1 -o docs -g cmd/server/main.go

RUN go build -o bin/server ./cmd/server/main.go

RUN touch /app/.env

CMD ["/app/bin/server"]
