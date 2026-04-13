# Advanced Programming 2 — Assignment 2

**gRPC Migration & Contract-First Development**

Student: Farida Dovletbayeva
Astana IT University
Course: Advanced Programming 2

---

# Overview

This project is a migration of a microservice system from REST to gRPC using a Contract-First approach.

The system consists of two services:

* Order Service (REST API + gRPC client)
* Payment Service (REST API + gRPC server)

---

# Architecture

* Clean Architecture is preserved
* Business logic (Use Cases) is unchanged from Assignment 1
* Only the communication layer is migrated from REST to gRPC

```
[ Client ]
    ↓ REST (Gin)
[ Order Service ]
    ↓ gRPC
[ Payment Service ]
    ↓
[ PostgreSQL ]
```

---

# Technologies Used

* Go (Golang)
* gRPC
* Protocol Buffers
* Gin (REST API)
* PostgreSQL
* Clean Architecture

---

# Repositories

### Proto Repository (Contract-First)

Contains `.proto` files:
https://github.com/Farida2025/assignment2-adp2

---

### Generated Code Repository

Contains generated `.pb.go` files:
https://github.com/Farida2025/assignment2-generated

---

# Protobuf Contract

```proto
syntax = "proto3";

package payment;

option go_package = "github.com/Farida2025/assignment2-generated/payment";

service PaymentService {
  rpc ProcessPayment (PaymentRequest) returns (PaymentResponse);
}

message PaymentRequest {
  string order_id = 1;
  int64 amount = 2;
}

message PaymentResponse {
  string transaction_id = 1;
  string status = 2;
}
```

---

# Configuration

### Environment Variables

| Variable          | Description                                       |
| ----------------- | ------------------------------------------------- |
| PAYMENT_GRPC_ADDR | Address of Payment Service (e.g. localhost:50051) |

---

# How to Run

## 1. Start PostgreSQL

Create databases:

```sql
CREATE DATABASE order_db;
CREATE DATABASE payment_db;
```

---

## 2. Run Payment Service

```bash
cd assignment1/payment-service
go run cmd/payment/main.go
```

Runs:

* gRPC server on :50051
* REST API on :8081

---

## 3. Run Order Service

### Windows PowerShell:

```powershell
$env:PAYMENT_GRPC_ADDR="localhost:50051"
go run cmd/order/main.go
```

### Linux / macOS:

```bash
PAYMENT_GRPC_ADDR=localhost:50051 go run cmd/order/main.go
```

Runs:

* REST API on :8080

---

# API Endpoints

## Order Service (REST)

| Method | Endpoint           |
| ------ | ------------------ |
| POST   | /orders            |
| GET    | /orders/:id        |
| PATCH  | /orders/:id/cancel |

---

## Payment Service (REST - optional)

| Method | Endpoint            |
| ------ | ------------------- |
| POST   | /payments           |
| GET    | /payments/:order_id |

---

# gRPC Communication

* Order Service acts as a gRPC client
* Payment Service acts as a gRPC server
* Communication uses strongly typed contracts

---

# Timeout Handling

Order Service uses a custom timeout (2 seconds):

```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
```

---

# Error Handling

Uses gRPC status codes:

* codes.Internal — server error
* codes.FailedPrecondition — payment declined

---

# Design Highlights

* Clean Architecture preserved
* Dependency Inversion applied
* No business logic in transport layer
* Contract-First development
* gRPC replaces REST between services
* Environment-based configuration (no hardcoding)

---

# Assignment Requirements Coverage

| Requirement           | Status |
| --------------------- | ------ |
| Contract-First        | Done   |
| gRPC Client/Server    | Done   |
| Clean Architecture    | Done   |
| Env Configuration     | Done   |
| Timeout               | Done   |
| Error Handling        | Done   |
| Repository Separation | Done   |

---

# Notes

* REST is kept only for external clients
* Internal communication is fully migrated to gRPC
* Generated code is imported via Go modules

---

# Conclusion

This project demonstrates a full migration from REST to gRPC while preserving architectural principles and improving type safety, performance, and maintainability.
