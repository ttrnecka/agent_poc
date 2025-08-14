package dto

type ProbeDTO struct {
	ID          string `json:"id"`
	CollectorID string `json:"collector_id"`
	Policy      string `json:"policy"`
	Version     string `json:"version"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	// TODO hide this in json when probe saveing is properly resolved
	Password string `json:"password"`
}
