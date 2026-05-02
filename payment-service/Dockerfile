FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o payment-app ./cmd/payment/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/payment-app .

EXPOSE 8081
EXPOSE 50051

CMD ["./payment-app"]