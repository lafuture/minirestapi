package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Name       string
	Pass       string
	DbName     string
	DbUrl      string
	Redis      string
	Kafka      string
	KafkaTopic string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env не найден")
	}

	Name = os.Getenv("NAME")
	Pass = os.Getenv("PASS")
	DbName = os.Getenv("DBNAME")
	DbUrl = os.Getenv("DBURL")
	Redis = os.Getenv("REDIS")
	Kafka = os.Getenv("KAFKA")
	KafkaTopic = os.Getenv("KAFKATOPIC")
}
