package services

import (
	"KafkaRedisTest/app/internal/domen/models"
	"context"
	"encoding/json"
	"fmt"
)

type RedisService struct {
	app *models.App
}

func NewRedisService(app *models.App) *RedisService {
	return &RedisService{
		app: app,
	}
}

func (s *RedisService) GetMessageFromCache(ctx context.Context, messageID string) (*models.Message, error) {
	value, err := s.app.RedisClient.Get(ctx, messageID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get message from redis: %w", err)
	}
	var message models.Message
	if err := json.Unmarshal([]byte(value), &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return &message, err
}
