package handler

import (
	"encoding/json"
	"myapp/internal/database"
	"myapp/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RegisterHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := db.GetUserByName(req.Name); err == nil {
			http.Error(w, "Пользователь уже существует", http.StatusBadRequest)
			return
		}

		err := db.RegisterUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Пользователь успешно зарегистрирован"))
	}
}

func LoginHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := db.GetUserByName(req.Name)
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
}

func GetPageHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		PageId := chi.URLParam(r, "id")
		IntPageId, err := strconv.Atoi(PageId)
		if err != nil {
			http.Error(w, "Некоректный номер", http.StatusBadRequest)
			return
		}

		p, err := db.GetPage(IntPageId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetPagesHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := db.GetPages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func MakePageHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.PageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := db.MakePage(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Страница успешно создана"))
	}
}

func EditPageHandler(db *database.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = db.EditPage(IntPageId, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Страница успешно отредактирована"))
	}
}
