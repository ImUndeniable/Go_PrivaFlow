package main

import (
	"log"
	"net/http"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres" // Update with your actual module path if different
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"
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

	// 4. Setup Router
	r := gin.Default()

	// 5. Create Endpoint
	r.POST("/request", func(c *gin.Context) {
		var req domain.ErasureRequest

		// BindJSON reads the user's input and fills the 'req' struct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to Database using our clean Repository
		if err := repo.Create(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Request received!", "data": req})
	})

	// 6. Start Server
	log.Println("🚀 PrivaFlow is live on :8080")
	r.Run(":8080")
}
