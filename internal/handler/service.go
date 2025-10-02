package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/database"
	"myapp/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	db    *database.Postgres
	kafka *kafka.Writer
	rd    *redis.Client
	us    *UserService
}

type UserService struct {
	db    *database.Postgres
	kafka *kafka.Writer
	rd    *redis.Client
}

func NewHandler(db *database.Postgres, kafka *kafka.Writer, rd *redis.Client) *Handler {
	return &Handler{db: db, kafka: kafka, rd: rd}
}
func (us *UserService) RegisterUser(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	if _, err := us.db.GetUserByName(req.Name); err == nil {
		return nil, fmt.Errorf("Пользователь уже существует")
	}

	u := models.User{
		Name:     req.Name,
		Password: req.Password,
		Mail:     req.Mail,
	}

	err := us.db.RegisterUser(u)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	if err := us.kafka.WriteMessages(ctx,
		kafka.Message{Value: data},
	); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *UserService) LoginUser(req models.LoginRequest) (*models.User, error) {
	u, err := us.db.GetUserByName(req.Name)
	if err != nil {
		return nil, err
	}

	if u.Password != req.Password {
		return nil, fmt.Errorf("Неправильно введен пароль")
	}

	return &u, nil
}

func (us *UserService) GetPageForUser(PageId int) (*models.Page, error) {
	p, err := us.db.GetPage(PageId)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (us *UserService) GetPagesForUser() (*[]models.Page, error) {
	p, err := us.db.GetPages()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (us *UserService) MakePageForUser(req models.PageRequest) error {
	err := us.db.MakePage(req)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) EditPageForUser(PageId int, req models.PageRequest) error {
	err := us.db.EditPage(PageId, req)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetLastPagesForUser() ([]byte, error) {
	cache, err := us.rd.Get(context.Background(), "last_pages").Result()
	if err == nil {
		return []byte(cache), err
	}

	p, err := us.db.GetPages()
	if err != nil {
		return nil, err
	}

	var lp []models.Page
	if len(p) > 10 {
		lp = p[len(p)-10:]
	} else {
		lp = p
	}

	data, err := json.Marshal(lp)
	if err != nil {
		return nil, err
	}

	if err := us.rd.Set(context.Background(), "last_pages", data, 600*time.Second).Err(); err != nil {
		return nil, err
	}

	return data, nil
}
