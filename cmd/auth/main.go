package main

import (
	"context"
	"log"
	"myapp/config"
	"myapp/internal/database"
	"myapp/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "users",
	})
	defer writer.Close()

	rd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rd.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db, err := database.New(config.DbUrl)
	if err != nil {
		log.Fatal("Failed to init storage", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handler.NewHandler(db, writer, rd)

	r.Post("/signup", h.Register)
	r.Post("/login", h.Login)
	r.Get("/page/{id}", h.GetPage)
	r.Get("/pages", h.GetPages)
	r.Post("/page", h.MakePage)
	r.Put("/page/{id}", h.EditPage)
	r.Get("/pages/last", h.GetLastPages)

	http.ListenAndServe(":3333", r)
}
