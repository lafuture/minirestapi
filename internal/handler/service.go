package handler

import (
	"myapp/internal/database"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	db    *database.Postgres
	kafka *kafka.Writer
	rd    *redis.Client
}

func NewHandler(db *database.Postgres, kafka *kafka.Writer, rd *redis.Client) *Handler {
	return &Handler{db: db, kafka: kafka, rd: rd}
}
