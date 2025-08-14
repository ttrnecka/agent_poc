package dto

type PolicyDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	FileName    string   `json:"file_name"`
	Versions    []string `json:"versions"`
}
