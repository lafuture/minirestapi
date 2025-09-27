package models

type RegisterRequest struct {
	Name     string `json:"text"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

type LoginRequest struct {
	Name     string `json:"text"`
	Password string `json:"password"`
}

type CheckRegisterRequest struct {
	Code int `json:"code"`
}

type PageRequest struct {
	Name string `json:"text"`
	Text string `json:"text"`
}
