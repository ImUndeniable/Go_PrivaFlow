package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/broker/kafka"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Setup Resources (Same as before)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	dsn := "host=" + dbHost + " user=admin password=secretpassword dbname=privaflow port=5432 sslmode=disable"
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	repo := postgres.NewErasureRepository(db)

	consumer := kafka.NewEventConsumer(kafkaBroker, "erasure-requests", "erasure-worker-group")
	defer consumer.Close() // Ensure connection closes on exit

	// 👇 2. NEW: Create a "Trap" for Ctrl+C
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to listen for OS signals (Syscall)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run a background listener
	go func() {
		<-sigChan // Block here until Ctrl+C is pressed
		log.Println("\n🛑 Shutdown signal received! Stopping new tasks...")
		cancel() // Cancel the context (This tells the loop to stop)
	}()

	log.Println("👷 Worker started! Waiting for tasks... (Press Ctrl+C to stop gracefully)")

	// 3. The Loop (Updated)
	for {
		// A. Check if we should quit BEFORE looking for work
		if ctx.Err() != nil {
			log.Println("🚪 Context closed, exiting loop.")
			break
		}

		// B. Fetch Message (Pass ctx so it unblocks if we quit while waiting)
		msg, err := consumer.FetchMessage(ctx)
		if err != nil {
			// If error is because we are shutting down, break the loop
			if ctx.Err() != nil {
				break
			}
			log.Println("⚠️ Error fetching:", err)
			continue
		}

		// --- CRITICAL SECTION START ---
		// Once we reach here, we MUST finish, even if Ctrl+C is pressed.

		var req domain.ErasureRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Println("⚠️ Invalid JSON:", err)
			continue
		}

		log.Printf("📥 Processing: %s (DO NOT STOP ME)", req.Email)

		// Simulate Work (5 seconds)
		// Try hitting Ctrl+C while this is running!
		time.Sleep(5 * time.Second)

		repo.UpdateStatus(req.Email, "COMPLETED")
		log.Printf("✅ Finished: %s", req.Email)

		// --- CRITICAL SECTION END ---
	}

	log.Println("👋 Worker exited cleanly. No data corruption!")
}
