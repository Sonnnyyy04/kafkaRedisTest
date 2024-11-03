package models

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	RedisAddr    string
	RedisDb      int
}
