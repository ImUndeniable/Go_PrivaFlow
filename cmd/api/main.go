package main

import (
	"log"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/handler"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres" // Update with your actual module path if different
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/services"
	"github.com/gin-gonic/gin"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=admin password=secretpassword dbname=privaflow port=5432 sslmode=disable"
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	// 2. Auto-Migrate (Create table if it doesn't exist)
	// In production, we use tools like 'golang-migrate', but this is fine for now.
	db.AutoMigrate(&domain.ErasureRequest{})

	// 3. Initialize Repository (The "Fridge")
	repo := postgres.NewErasureRepository(db)
	// The Chef (uses the Fridge)
	svc := services.NewErasureService(repo)
	// The Waiter (uses chef)
	h := handler.NewErasureHandler(svc)

	// 4. Setup Router
	r := gin.Default()

	// 5. Route Definition
	r.POST("/request", h.RequestErasure)

	// 6. Start Server
	log.Println("🚀 PrivaFlow is live on :8080")
	r.Run(":8080")
}
