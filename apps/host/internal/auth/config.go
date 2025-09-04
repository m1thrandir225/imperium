package auth

import "time"

type Config struct {
	User                  User   `json:"user" mapstructure:"user"`
	AccessToken           string `json:"access_token" mapstructure:"access_token"`
	RefreshToken          string `json:"refresh_token" mapstructure:"refresh_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at" mapstructure:"access_token_expires_at"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at" mapstructure:"refresh_token_expires_at"`
}

func (c *Config) GetCurrentUser() User {
	return c.User
}

func (c *Config) GetAccessToken() string {
	return c.AccessToken
}

func (c *Config) GetRefreshToken() string {
	return c.RefreshToken
}

func (c *Config) GetAccessTokenExpiresAt() time.Time {
	if c.AccessTokenExpiresAt == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse(time.RFC3339, c.AccessTokenExpiresAt)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}

func (c *Config) GetRefreshTokenExpiresAt() time.Time {
	if c.RefreshTokenExpiresAt == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse(time.RFC3339, c.RefreshTokenExpiresAt)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}

func (c *Config) IsAccessTokenExpired() bool {
	return time.Now().After(c.GetAccessTokenExpiresAt())
}

func (c *Config) IsRefreshTokenExpired() bool {
	return time.Now().After(c.GetRefreshTokenExpiresAt())
}

func (c *Config) SetCurrentUser(user User) {
	c.User = user
}

func (c *Config) SetAccessToken(accessToken string) {
	c.AccessToken = accessToken
}

func (c *Config) SetRefreshToken(refreshToken string) {
	c.RefreshToken = refreshToken
}

func (c *Config) SetAccessTokenExpiresAt(accessTokenExpiresAt time.Time) {
	c.AccessTokenExpiresAt = accessTokenExpiresAt.Format(time.RFC3339)
}

func (c *Config) SetRefreshTokenExpiresAt(refreshTokenExpiresAt time.Time) {
	c.RefreshTokenExpiresAt = refreshTokenExpiresAt.Format(time.RFC3339)
}

func LoadConfig(user User, accessToken string, refreshToken string, accessTokenExpiresAt time.Time, refreshTokenExpiresAt time.Time) *Config {
	return &Config{
		User:                  user,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt.Format(time.RFC3339),
		RefreshTokenExpiresAt: refreshTokenExpiresAt.Format(time.RFC3339),
	}
}

// NewDefaultConfig returns a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		User:                  User{},
		AccessToken:           "",
		RefreshToken:          "",
		AccessTokenExpiresAt:  "",
		RefreshTokenExpiresAt: "",
	}
}
