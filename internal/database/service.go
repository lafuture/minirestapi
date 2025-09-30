package database

import (
	"context"
	"log"
	"myapp/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
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
		config.Migra,
		Url,
	)
	if err != nil {
		log.Fatal("Не удалось инициализировать миграции:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Не удалось применить миграции:", err)
	} else {
		log.Println("Миграции применены или изменений не было")
	}

	log.Println("Миграции успешно применены")

	return &Postgres{db: db}, nil
}
