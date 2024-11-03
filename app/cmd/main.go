package main

import (
	"KafkaRedisTest/app/internal/domen/models"
	"KafkaRedisTest/app/internal/domen/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	config := models.Config{
		KafkaBrokers: []string{"kafka:9092"},
		KafkaTopic:   "test-topic",
		RedisAddr:    "redis:6379",
		RedisDb:      0,
	}
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.KafkaBrokers, kafkaConfig)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer producer.Close()
	consumer, err := sarama.NewConsumer(config.KafkaBrokers, kafkaConfig)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer consumer.Close()
	redisClient := redis.NewClient(&redis.Options{Addr: config.RedisAddr, DB: config.RedisDb})
	app := &models.App{
		Config:        config,
		KafkaProducer: producer,
		KafkaConsumer: consumer,
		RedisClient:   redisClient,
	}
	kafkaService := services.NewKafkaService(app)
	redisService := services.NewRedisService(app)
	ctx := context.Background()
	message := models.Message{
		ID:        "123",
		Content:   "Hello, Kafka & Redis!",
		Timestamp: time.Now().UTC(),
	}
	if err := kafkaService.PublishMessage(ctx, message); err != nil {
		log.Fatalf("failed to publish message: %v", err)
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
	}
	if err := redisClient.Set(ctx, message.ID, msgBytes, 0).Err(); err != nil {
		log.Fatalf("failed to set message in redis: %v", err)
	}
	cachedMessage, err := redisService.GetMessageFromCache(ctx, message.ID)
	if err != nil {
		log.Fatalf("failed to get message from cache: %v", err)
	}
	fmt.Printf("Message from cache: %+v\n", cachedMessage)
}
