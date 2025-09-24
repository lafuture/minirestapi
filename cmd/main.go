package main

import (
	"log"
	"myapp/config"
	"myapp/internal/database"
	"myapp/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.New(config.DbUrl)
	if err != nil {
		log.Fatal("Failed to init storage", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/signup", handler.RegisterHandler(db))
	r.Post("/login", handler.LoginHandler(db))
	r.Get("/page/{id}", handler.GetPageHandler(db))
	r.Get("/pages", handler.GetPagesHandler(db))
	r.Post("/page", handler.MakePageHandler(db))
	r.Put("/page/{id}", handler.EditPageHandler(db))

	http.ListenAndServe(":3333", r)
}
