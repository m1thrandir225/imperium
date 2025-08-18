package internal

import "time"

type ClientDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type HostDTO struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
}

type UserDTO struct {
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Client    []ClientDTO `json:"client"`
	Hosts     []HostDTO   `json:"hosts"`
}
