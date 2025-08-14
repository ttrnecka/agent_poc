package dto

type CollectorDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Password string `json:"password"`
}
