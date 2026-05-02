# Advanced Programming 2 — Assignment 3: Event-Driven Architecture

## Overview
This project evolves the microservice system from Assignment 2 by introducing **Event-Driven Architecture (EDA)**.  
The system now uses **RabbitMQ** for asynchronous communication between services to handle notifications and ensure reliability.

---

## Architecture & Event Flow

The system follows a **Choreography-based EDA** pattern.  
Communication between **Order** and **Payment** remains synchronous (gRPC), while **Payment** and **Notification** interact asynchronously via RabbitMQ.

### Flow of Events:
1. **Order Service** receives a REST request and calls **Payment Service** via gRPC.
2. **Payment Service** processes the payment, stores it in `payment_db`, and publishes a `payment.completed` event to RabbitMQ.
3. **RabbitMQ Broker** routes the message to the queue or to a **Dead Letter Queue (DLQ)** if processing fails.
4. **Notification Service** consumes events, checks for duplicates (Idempotency), and sends a notification.

---

## 🛠 Engineering Decisions

### 1. Idempotency Strategy
To prevent duplicate notifications (e.g., message redelivery), an **in-memory idempotency mechanism** is used.

- **Mechanism**: `map[string]bool` with `sync.Mutex`
- **Key**: `OrderID`
- **Logic**:
    - If `OrderID` already processed → skip + ACK
    - Otherwise → process and store


---

### 2. Manual ACK & Reliability

- **Auto-ACK disabled**
- Messages are acknowledged **only after successful processing**

#### Behavior:
-  Success → `d.Ack(false)`
-  Failure → `d.Nack(false, true)` (requeue)

This guarantees:

> **At-Least-Once Delivery**

If the service crashes during processing:
- message stays in queue
- gets redelivered

---

### 3. Dead Letter Queue (DLQ)

To handle invalid or permanently failing messages:

- **Dead Letter Exchange**: `payment.dlx`
- **Dead Letter Queue**: `payment.failed`

#### Logic:
- If message is invalid → `d.Nack(false, false)`
- RabbitMQ moves it to DLQ

#### Example:
- Simulated failure for amount = `666.66`

#### Result:
- avoids infinite retries
- enables debugging via RabbitMQ UI

---

##  How to Run with Docker

### 1. Build and Start

```bash
docker-compose up --build