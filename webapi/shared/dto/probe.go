package dto

type ProbeDTO struct {
	ID          string `json:"id"`
	CollectorID string `json:"collector_id" validate:"required"`
	Policy      string `json:"policy" validate:"required"`
	Version     string `json:"version" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Port        int    `json:"port" validate:"required"`
	User        string `json:"user" validate:"required"`
	// TODO hide this in json when probe saveing is properly resolved
	Password string `json:"password"`
}
