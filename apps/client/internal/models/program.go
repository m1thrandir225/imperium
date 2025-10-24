package models

type Program struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	HostID      string `json:"hostId"`
}
