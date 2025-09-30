package database

import (
	"context"
	"fmt"
	"myapp/internal/models"
)

func (p *Postgres) RegisterUser(m models.User) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO users (name, password, mail) VALUES ($1, $2, $3)",
		m.Name, m.Password, m.Mail)
	if err != nil {
		return fmt.Errorf("Не удалось зарегистрировать пользователя: %w", err)
	}

	return nil
}

func (p *Postgres) GetUserByName(name string) (models.User, error) {
	var user models.User
	err := p.db.QueryRow(context.Background(),
		"SELECT * FROM users WHERE name = $1", name).
		Scan(&user.Id, &user.Name, &user.Password, &user.Mail)
	if err != nil {
		return models.User{}, fmt.Errorf("Не удалось найти пользователя: %w", err)
	}

	return user, nil
}

func (p *Postgres) MakePage(m models.PageRequest) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO pages (name, text) VALUES ($1, $2)",
		m.Name, m.Text)
	if err != nil {
		return fmt.Errorf("Не удалось сделать статью: %w", err)
	}

	return nil
}

func (p *Postgres) EditPage(id int, m models.PageRequest) error {
	_, err := p.db.Exec(context.Background(),
		"UPDATE pages SET name = $1, text = $2 WHERE id = $3",
		m.Name, m.Text, id)
	if err != nil {
		return fmt.Errorf("Не удалось отредактировать статью: %w", err)
	}

	return nil
}

func (p *Postgres) GetPage(id int) (models.Page, error) {
	var page models.Page
	err := p.db.QueryRow(context.Background(),
		"SELECT id, name, text FROM pages WHERE id = $1", id).
		Scan(&page.Id, &page.Name, &page.Text)
	if err != nil {
		return models.Page{}, fmt.Errorf("Не удалось посмотреть статью: %w", err)
	}

	return page, nil
}

func (p *Postgres) GetPages() ([]models.Page, error) {
	rows, err := p.db.Query(context.Background(), "SELECT id, name, text FROM pages")
	if err != nil {
		return nil, fmt.Errorf("Не удалось посмотреть все статьи: %w", err)
	}
	defer rows.Close()

	pages := []models.Page{}
	for rows.Next() {
		var page models.Page
		if err := rows.Scan(&page.Id, &page.Name, &page.Text); err != nil {
			return nil, fmt.Errorf("Ошибка чтения строки: %w", err)
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (p *Postgres) Close() {
	p.db.Close(context.Background())
}
