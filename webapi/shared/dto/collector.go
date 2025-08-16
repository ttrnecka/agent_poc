package dto

type CollectorDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Status   string `json:"status"`
	Password string `json:"password"`
}
