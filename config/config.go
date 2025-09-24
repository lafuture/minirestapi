package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Name   string
	Pass   string
	DbName string
	DbUrl  string
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
}
