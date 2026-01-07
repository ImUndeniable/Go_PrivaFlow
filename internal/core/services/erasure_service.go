// CHEF
package services

import (
	"context"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/broker/kafka"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/adapter/storage/postgres"
	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/domain"
)

type ErasureService struct {
	repo     *postgres.ErasureRepository
	producer *kafka.EventProducer
}

func NewErasureService(repo *postgres.ErasureRepository, producer *kafka.EventProducer) *ErasureService {
	return &ErasureService{repo: repo, producer: producer}
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
	// 2. Publish to Kafka (The Shout) 📣
	// We use the Email as the "Key" so all requests for the same user go to the same partition
	err := r.producer.Publish(context.Background(), req.Email, req)
	if err != nil {
		// In a real system, we might retry or log a serious error here
		return nil, err
	}

	return req, nil
}
