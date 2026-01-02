package domain

import "time"

// ErasureRequest represents the data structure for a user's request to be forgotten.
// We add JSON tags (for the API) and GORM tags (for the Database).

type ErasureRequest struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null"`
	Status    string    `gorm:"default:'PENDING'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
