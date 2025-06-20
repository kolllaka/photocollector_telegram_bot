FROM golang:1.23-alpine AS builder

WORKDIR /tg-bot

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o ./start ./cmd/main.go



FROM alpine AS runner

WORKDIR /app

COPY --from=builder /tg-bot/start ./start
COPY .env .env
COPY locale.yaml locale.yaml

CMD ["/app/start"]
