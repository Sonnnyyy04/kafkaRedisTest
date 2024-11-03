package services

import (
	"KafkaRedisTest/app/internal/domen/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaService struct {
	app *models.App
}

func NewKafkaService(app *models.App) *KafkaService {
	return &KafkaService{
		app: app,
	}
}

func (s *KafkaService) PublishMessage(ctx context.Context, msg models.Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	partition, offset, err := s.app.KafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: s.app.Config.KafkaTopic,
		Value: sarama.StringEncoder(msgBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}
	log.Printf("Message published to partition %d at offset %d", partition, offset)
	return nil
}
