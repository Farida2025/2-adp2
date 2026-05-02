FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o notification-app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/notification-app .

CMD ["./notification-app"]