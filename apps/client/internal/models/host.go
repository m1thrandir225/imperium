package models

type Host struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
}

type SimpleHost struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
