package models

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CheckRegisterRequest struct {
	Code int `json:"code"`
}

type PageRequest struct {
	Name string `json:"name"`
	Text string `json:"text"`
}
