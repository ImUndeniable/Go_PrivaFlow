package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/broker/kafka"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup DB
	dsn := "host=localhost user=admin password=secretpassword dbname=privaflow port=5432 sslmode=disable"
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Worker failed to connect to DB:", err)
	}
	repo := postgres.NewErasureRepository(db)

	// 2. Setup Kafka Consumer
	consumer := kafka.NewEventConsumer("localhost:9092", "erasure-requests", "erasure-worker-group")
	defer consumer.Close()

	log.Println("👷 Worker started! Waiting for tasks...")

	// 3. Infinite Loop
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Println("⚠️ Error fetching:", err)
			continue
		}

		var req domain.ErasureRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Println("⚠️ Invalid JSON:", err)
			continue
		}

		log.Printf("📥 Processing: %s", req.Email)

		// Simulate Work
		time.Sleep(5 * time.Second)

		// Update DB
		repo.UpdateStatus(req.Email, "COMPLETED")
		log.Printf("✅ Finished: %s", req.Email)
	}
}
