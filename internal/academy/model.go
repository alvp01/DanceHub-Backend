// internal/academy/model.go
package academy

import "time"

type Academy struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PrimaryPhone   string    `json:"primary_phone"`
	SecondaryPhone *string   `json:"secondary_phone,omitempty"` // puntero = nullable
	PasswordHash   string    `json:"-"`                         // nunca se serializa
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RegisterRequest es lo que recibe el endpoint
type RegisterRequest struct {
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	PrimaryPhone   string  `json:"primary_phone"`
	SecondaryPhone *string `json:"secondary_phone"`
	Password       string  `json:"password"`
}

// RegisterResponse es lo que devuelve el endpoint (sin hash)
type RegisterResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PrimaryPhone   string    `json:"primary_phone"`
	SecondaryPhone *string   `json:"secondary_phone,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// internal/academy/model.go (agregar al archivo existente)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // "Bearer"
	ExpiresIn    int    `json:"expires_in"` // segundos
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID        string     `json:"id"`
	AcademyID string     `json:"academy_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}
