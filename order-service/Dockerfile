FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o order-app ./cmd/order/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/order-app .

EXPOSE 8080

CMD ["./order-app"]