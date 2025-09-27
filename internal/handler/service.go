package handler

import "myapp/internal/database"

type Handler struct {
	db *database.Postgres
}

func NewHandler(db *database.Postgres) *Handler {
	return &Handler{db: db}
}
