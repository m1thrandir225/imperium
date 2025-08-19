package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateHostRequest struct {
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}
