// CHEF
package services

import (
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"
)

type ErasureService struct {
	repo *postgres.ErasureRepository
}

func NewErasureService(repo *postgres.ErasureRepository) *ErasureService {
	return &ErasureService{repo: repo}
}

func (r *ErasureService) RequestErasure(email string) (*domain.ErasureRequest, error) {
	// 1. Chef prepares the ingredients (Business Logic)
	req := &domain.ErasureRequest{
		Email:  email,
		Status: "PENDING",
	}
	// 2. Chef opens the Fridge (The Next Link)
	// 👇 THIS MOVES IT TO THE DB CODE 👇
	if err := r.repo.Create(req); err != nil {
		return nil, err
	}

	return req, nil
}
