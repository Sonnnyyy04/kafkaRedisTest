package models

import (
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config        Config
	KafkaProducer sarama.SyncProducer
	KafkaConsumer sarama.Consumer
	RedisClient   *redis.Client
}
