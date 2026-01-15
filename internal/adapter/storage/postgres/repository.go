package postgres

import (
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain" // Make sure this path matches your go.mod name!
	"gorm.io/gorm"
)

type ErasureRepository struct {
	db *gorm.DB
}

func NewErasureRepository(db *gorm.DB) *ErasureRepository {
	return &ErasureRepository{db: db}
}

// Create is a method attached to the struct
func (r *ErasureRepository) Create(req *domain.ErasureRequest) error {
	// 👇 THIS IS THE FINAL DESTINATION 👇
	return r.db.Create(req).Error
}

// UpdateStatus finds a request by Email and updates its status
func (r *ErasureRepository) UpdateStatus(email string, status string) error {
	return r.db.Model(&domain.ErasureRequest{}).
		Where("email = ?", email).
		Update("status", status).Error
}
