package models

type Response struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Output  string `json:"output"`
}
