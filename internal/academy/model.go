// internal/academy/model.go
package academy

import "time"

type Academy struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PrimaryPhone   string    `json:"primary_phone"`
	SecondaryPhone *string   `json:"secondary_phone,omitempty"`
	PasswordHash   string    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
type RegisterRequest struct {
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	PrimaryPhone   string  `json:"primary_phone"`
	SecondaryPhone *string `json:"secondary_phone"`
	Password       string  `json:"password"`
}

type RegisterResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	PrimaryPhone   string  `json:"primary_phone"`
	SecondaryPhone *string `json:"secondary_phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID        string     `json:"id"`
	AcademyID string     `json:"academy_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"-"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}
