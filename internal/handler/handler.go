package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"myapp/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

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

	random := rand.Intn(9000) + 1000
	fmt.Println(random)

	session, _ := store.Get(r, "signup-session")
	session.Values["mail"] = req.Mail
	session.Values["name"] = req.Name
	session.Values["password"] = req.Password
	session.Values["code"] = random
	session.Save(r, w)

	http.Redirect(w, r, "/signup/check", http.StatusSeeOther)
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

func (h *Handler) CheckRegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Введите код который пришел на почту"))
}

func (h *Handler) CheckRegister(w http.ResponseWriter, r *http.Request) {
	var req models.CheckRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, "signup-session")

	if int(req.Code) != session.Values["code"].(int) {
		http.Error(w, "Неправильно введен код", http.StatusBadRequest)
		return
	}

	u := models.User{
		Name:     session.Values["name"].(string),
		Password: session.Values["password"].(string),
		Mail:     session.Values["mail"].(string),
	}

	err := h.db.RegisterUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно создан"))
}
