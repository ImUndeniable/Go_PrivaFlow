# 🛡️ PrivaFlow - Distributed Privacy Request System

PrivaFlow is a scalable, event-driven backend system built in **Go (Golang)** designed to handle high-throughput user privacy requests (e.g., "Right to Erasure"). It uses a microservices approach with **Kafka** for asynchronous processing and **Docker** for containerization.

## 🏗️ Architecture (The "Nano Banana" Flow)

**User** → **API (Producer)** → **Kafka (Broker)** → **Worker (Consumer)** → **DB**

1.  **API Service:** Accepts HTTP requests, saves initial status ("PENDING") to Postgres, and publishes an event to Kafka. Returns `200 OK` immediately.
2.  **Kafka:** Buffers the message, ensuring no data is lost even if the worker is busy.
3.  **Worker Service:** Listens to Kafka, picks up the task, simulates heavy processing (AI Scanning), and updates the DB status to "COMPLETED".

## 🛠️ Tech Stack
* **Language:** Go (Golang) 1.25+
* **Framework:** Gin (HTTP Web Framework)
* **Database:** PostgreSQL (with GORM ORM)
* **Message Broker:** Apache Kafka & Zookeeper
* **Containerization:** Docker & Docker Compose
* **Architecture:** Clean Architecture + Event-Driven

## 🚀 How to Run (One Command)
Prerequisite: Docker Desktop installed.

```bash
# Clone the repo
git clone [https://github.com/YOUR_USERNAME/Go_PrivaFlow.git](https://github.com/YOUR_USERNAME/Go_PrivaFlow.git)
cd Go_PrivaFlow

# Start the entire system
docker-compose up --build
