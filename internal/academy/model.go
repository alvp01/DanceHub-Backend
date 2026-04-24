// internal/academy/model.go
package academy

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Academy struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name           string    `json:"name" gorm:"column:name;not null;uniqueIndex"`
	Email          string    `json:"email" gorm:"column:email;not null;uniqueIndex"`
	PrimaryPhone   string    `json:"primary_phone" gorm:"column:primary_phone;not null;uniqueIndex"`
	SecondaryPhone *string   `json:"secondary_phone,omitempty" gorm:"column:secondary_phone;uniqueIndex"`
	PasswordHash   string    `json:"-" gorm:"column:password_hash;not null"`
	CreatedAt      time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"-" gorm:"column:updated_at"`
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
	ID        string     `json:"id" gorm:"type:uuid;primaryKey"`
	AcademyID string     `json:"academy_id" gorm:"column:academy_id;type:uuid;not null;index"`
	TokenHash string     `json:"-" gorm:"column:token_hash;not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"column:expires_at;not null"`
	CreatedAt time.Time  `json:"-" gorm:"column:created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"column:revoked_at"`
}

func (a *Academy) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}

func (rt *RefreshToken) BeforeCreate(_ *gorm.DB) error {
	if rt.ID == "" {
		rt.ID = uuid.NewString()
	}
	return nil
}
