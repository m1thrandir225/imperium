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
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Name      string `json:"name"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"token"`
}
