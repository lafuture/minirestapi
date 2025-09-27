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

	h := handler.NewHandler(db)

	r.Post("/signup", h.Register)
	r.Post("/login", h.Login)
	r.Get("/page/{id}", h.GetPage)
	r.Get("/pages", h.GetPages)
	r.Post("/page", h.MakePage)
	r.Put("/page/{id}", h.EditPage)
	r.Get("/signup/check", h.CheckRegisterPage)
	r.Post("/signup/check", h.CheckRegister)

	http.ListenAndServe(":3333", r)
}
