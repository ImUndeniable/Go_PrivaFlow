package main

import (
	"log"
	"os"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/broker/kafka"
	http "github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/handler"
	middleware "github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/handler/http"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres" // Update with your actual module path if different
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/services"
	"github.com/gin-gonic/gin"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // Default for local run
	}

	dsn := "host=" + dbHost + " user=admin password=secretpassword dbname=privaflow port=5432 sslmode=disable"
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092" // Default
	}
	// 2. Auto-Migrate (Create table if it doesn't exist)
	// In production, we use tools like 'golang-migrate', but this is fine for now.
	db.AutoMigrate(&domain.ErasureRequest{})

	producer := kafka.NewEventProducer(kafkaBroker, "erasure-requests")
	defer producer.Close() // Close connection when app stops

	// 3. Initialize Repository (The "Fridge")
	repo := postgres.NewErasureRepository(db)
	// The Chef (uses the Fridge)
	svc := services.NewErasureService(repo, producer)
	// The Waiter (uses chef)
	//h := handler.NewErasureHandler(svc)
	erasureHandler := http.NewErasureHandler(svc)

	// 4. Setup Router
	r := gin.Default()

	// 5. Route Definition
	//r.POST("/request", h.RequestErasure)

	r.POST("/login", func(c *gin.Context) {
		token, err := middleware.GenerateToken("admin@privaflow.com")
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(200, gin.H{
			"token":   token,
			"message": "Here is your badge! Use it in the Authorization header.",
		})
	})

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Move your existing handler inside this group
		protected.POST("/request", erasureHandler.RequestErasure)
	}

	// 6. Start Server
	log.Println("🚀 PrivaFlow is live on :8081")
	r.Run(":8081")
}
