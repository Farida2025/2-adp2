# Assignment 1 – Clean Architecture based Microservices (Order & Payment)

**Student:** Farida Dovletbaeva  
**Course:** Advanced Programming 2  
**Date:** April 2026

## Overview

This project implements a simple two-service platform consisting of **Order Service** and **Payment Service**.  
Both services follow **Clean Architecture** principles and demonstrate basic **Microservices** concepts:
- Bounded Contexts
- Database per Service
- Synchronous REST communication with timeout
- Separation of Concerns and Dependency Inversion


Each service has its own:
- Domain models
- Use cases (business logic)
- Repository layer (with interface)
- HTTP handlers (thin layer)
- Separate PostgreSQL database

## Technologies

- Go 1.23
- Gin (HTTP framework)
- PostgreSQL (Database per Service)
- Clean Architecture
- REST communication

## How to Run

1. Make sure PostgreSQL is running and two databases are created:
    - `order_db`
    - `payment_db`

2. Apply migrations (run SQL scripts from `migrations/schema.sql` in each database).

3. Start **Payment Service** first:
   ```bash
   cd payment-service
   go run cmd/payment/main.go

## Architecture Decisions

- **Clean Architecture**:
    - Business logic is located in the Use Cases layer
    - Handlers are thin — they only parse requests and return responses
    - Use cases depend on interfaces (ports), not on concrete implementations
    - Composition root (dependency wiring) is done in `main.go`

- **Microservices**:
    - Separate databases for each service (`order_db` and `payment_db`)
    - No shared models or code between the two services
    - Order Service communicates with Payment Service via REST with a 2-second timeout

- **Resilience**:
    - If Payment Service is down or slow, Order Service returns HTTP 503 Service Unavailable
    - The order status is changed to "Failed"

- **Business Rules**:
    - Money amount is stored as `int64` (in cents) for financial accuracy
    - If amount > 100000, Payment Service returns "Declined"
    - Only orders with status "Pending" can be cancelled

## Failure Scenario

When the Payment Service is unavailable:

- `http.Client` timeout of 2 seconds is triggered
- Order status is updated to "Failed"
- The API returns HTTP 503 Service Unavailable

## Diagram

See `diagram.png` in the root folder.

The diagram illustrates:
- Two independent services with Clean Architecture layers (Domain → Use Case → Repository → Transport)
- REST communication between Order Service and Payment Service
- Separate databases for each service

## Conclusion

This project demonstrates the proper application of **Clean Architecture** principles inside each service and basic **Microservices** concepts (bounded contexts, database per service, and resilient synchronous communication) as taught in Lecture 1 and Lecture 2.
