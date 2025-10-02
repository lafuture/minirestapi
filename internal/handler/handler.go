package handler

import (
	"encoding/json"
	"fmt"
	"myapp/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.us.RegisterUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Пользователь %s успешно создан", user.Name)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.us.LoginUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Пользователь %s успешно залогинился", user.Name)))
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	PageId := chi.URLParam(r, "id")
	IntPageId, err := strconv.Atoi(PageId)
	if err != nil {
		http.Error(w, "Некоректный номер", http.StatusBadRequest)
		return
	}

	p, err := h.us.GetPageForUser(IntPageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetPages(w http.ResponseWriter, r *http.Request) {
	p, err := h.us.GetPagesForUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	err := h.us.MakePageForUser(req)
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

	if err := h.us.EditPageForUser(IntPageId, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Страница успешно отредактирована"))
}

func (h *Handler) GetLastPages(w http.ResponseWriter, r *http.Request) {
	data, err := h.us.GetLastPagesForUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
