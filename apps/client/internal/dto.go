package internal

import "time"

type ClientDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type SimpleHostDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HostDTO struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	Name      string `json:"name"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
}

type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
