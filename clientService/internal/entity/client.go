package entity

import "time"


type Client struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Password     string `json:"password"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}