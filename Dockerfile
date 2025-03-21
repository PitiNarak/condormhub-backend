FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/server/main.go

RUN touch /app/.env

CMD ["/app/bin/server"]