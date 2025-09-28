package handler

import (
	"context"
	"encoding/json"
	"myapp/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.db.GetUserByName(req.Name); err == nil {
		http.Error(w, "Пользователь уже существует", http.StatusBadRequest)
		return
	}

	u := models.User{
		Name:     req.Name,
		Password: req.Password,
		Mail:     req.Mail,
	}

	err := h.db.RegisterUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.kafka.WriteMessages(r.Context(),
		kafka.Message{Value: data},
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно создан"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.db.GetUserByName(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Password != req.Password {
		http.Error(w, "Неправильно введен пароль", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно залогинился"))
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	PageId := chi.URLParam(r, "id")
	IntPageId, err := strconv.Atoi(PageId)
	if err != nil {
		http.Error(w, "Некоректный номер", http.StatusBadRequest)
		return
	}

	p, err := h.db.GetPage(IntPageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetPages(w http.ResponseWriter, r *http.Request) {
	p, err := h.db.GetPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) MakePage(w http.ResponseWriter, r *http.Request) {
	var req models.PageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.MakePage(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Страница успешно создана"))
}

func (h *Handler) EditPage(w http.ResponseWriter, r *http.Request) {
	PageId := chi.URLParam(r, "id")
	IntPageId, err := strconv.Atoi(PageId)
	if err != nil {
		http.Error(w, "Некоректный номер", http.StatusBadRequest)
		return
	}

	var req models.PageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.EditPage(IntPageId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Страница успешно отредактирована"))
}

func (h *Handler) GetLastPages(w http.ResponseWriter, r *http.Request) {
	cache, err := h.rd.Get(context.Background(), "last_pages").Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cache))
		return
	}

	p, err := h.db.GetPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var lp []models.Page
	if len(p) > 10 {
		lp = p[len(p)-10:]
	} else {
		lp = p
	}

	data, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.rd.Set(context.Background(), "last_pages", data, 600*time.Second).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
