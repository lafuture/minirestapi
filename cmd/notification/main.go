package main

import (
	"context"
	"encoding/json"
	"log"
	"myapp/config"
	"myapp/internal/models"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Kafka},
		Topic:   config.KafkaTopic,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Ошибка при получении:", err)
		}

		var u models.User
		if err := json.Unmarshal(msg.Value, &u); err != nil {
			log.Fatal("Ошибка при распаковке:", err)
		}

		log.Printf("Письмо пользователю %s отправлено на %s\n", u.Name, u.Mail)
	}

}
