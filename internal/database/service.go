package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"log"
)

type Postgres struct {
	db *pgx.Conn
}

func New(Url string) (*Postgres, error) {
	db, err := pgx.Connect(context.Background(), Url)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	}

	m, err := migrate.New(
		"file://./migrations",
		Url,
	)
	if err != nil {
		log.Fatal("Не удалось инициализировать миграции:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Не удалось применить миграции:", err)
	}

	log.Println("Миграции успешно применены")

	return &Postgres{db: db}, nil
}
