package auth

import "time"

type User struct {
	ID        string `json:"id" mapstructure:"id"`
	Email     string `json:"email" mapstructure:"email"`
	Name      string `json:"name" mapstructure:"name"`
	CreatedAt string `json:"created_at" mapstructure:"created_at"`
	UpdatedAt string `json:"updated_at" mapstructure:"updated_at"`
}

func (u *User) GetCreatedAt() time.Time {
	if u.CreatedAt == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse(time.RFC3339, u.CreatedAt)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}

func (u *User) GetUpdatedAt() time.Time {
	if u.UpdatedAt == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse(time.RFC3339, u.UpdatedAt)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
