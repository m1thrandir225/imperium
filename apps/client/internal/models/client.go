package models

type Client struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

type RegisterClientRequest struct {
	Name      string `json:"client_name" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
}
type RegisterClientResponse struct {
	ID        string `json:"id"`
	Name      string `json:"client_name"`
	IPAddress string `json:"ip_address"`
}
