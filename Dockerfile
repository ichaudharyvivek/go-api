FROM golang:alpine

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest

RUN go build -o ./bin/api ./cmd/api

RUN go build -o ./bin/migrate ./cmd/migrate

CMD ["air", "-c", ".air.toml"]

EXPOSE 8080