// internal/academy/repository.go
package academy

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AcademyError string

func (e AcademyError) Error() string {
	return string(e)
}

const (
	ErrEmailAlreadyExists   AcademyError = "el email ya está registrado"
	ErrNameAlreadyExists    AcademyError = "el nombre ya está registrado"
	ErrPhoneAlreadyExists   AcademyError = "el teléfono ya está registrado"
	ErrAcademyNotFound      AcademyError = "academia no encontrada"
	ErrInvalidCredentials   AcademyError = "credenciales inválidas"
	ErrRefreshTokenNotFound AcademyError = "refresh token no encontrado"
	ErrRefreshTokenRevoked  AcademyError = "refresh token revocado"
	ErrRefreshTokenExpired  AcademyError = "refresh token expirado"
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *Academy) error {
	err := r.db.WithContext(ctx).Create(a).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "academies_email_key", "idx_academies_email":
				return ErrEmailAlreadyExists
			case "academies_name_key", "idx_academies_name":
				return ErrNameAlreadyExists
			case "academies_primary_phone_key", "idx_academies_primary_phone":
				return ErrPhoneAlreadyExists
			}
		}
		// Fallback: check error message for unique constraint on email
		if errMsg := err.Error(); errMsg != "" {
			if contains(errMsg, "UNIQUE constraint failed: academies.email") || contains(errMsg, "duplicate key value violates unique constraint 'academies_email_key'") || contains(errMsg, "duplicate key value violates unique constraint 'idx_academies_email'") {
				return ErrEmailAlreadyExists
			}
			if contains(errMsg, "UNIQUE constraint failed: academies.name") || contains(errMsg, "duplicate key value violates unique constraint 'academies_name_key'") || contains(errMsg, "duplicate key value violates unique constraint 'idx_academies_name'") {
				return ErrNameAlreadyExists
			}
			if contains(errMsg, "UNIQUE constraint failed: academies.primary_phone") || contains(errMsg, "duplicate key value violates unique constraint 'academies_primary_phone_key'") || contains(errMsg, "duplicate key value violates unique constraint 'idx_academies_primary_phone'") {
				return ErrPhoneAlreadyExists
			}
		}
		return fmt.Errorf("repository.Create: %w", err)
	}
	return nil
}

// contains is a helper for substring search (for error message fallback)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (len(s) > len(substr) && (s[0:len(substr)] == substr || contains(s[1:], substr)))))
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*Academy, error) {
	a := &Academy{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(a).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrAcademyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindByEmail: %w", err)
	}

	return a, nil
}

func (r *Repository) SaveRefreshToken(
	ctx context.Context,
	academyID string,
	rawToken string,
	expiresAt time.Time,
) error {
	err := r.db.WithContext(ctx).Create(&RefreshToken{
		AcademyID: academyID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: expiresAt,
	}).Error
	if err != nil {
		return fmt.Errorf("repository.SaveRefreshToken: %w", err)
	}
	return nil
}

func (r *Repository) FindRefreshToken(ctx context.Context, rawToken string) (*RefreshToken, error) {
	rt := &RefreshToken{}
	err := r.db.WithContext(ctx).Where("token_hash = ?", hashToken(rawToken)).First(rt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindRefreshToken: %w", err)
	}

	// Verificar estado en la capa de datos
	if rt.RevokedAt != nil {
		return nil, ErrRefreshTokenRevoked
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	return rt, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, rawToken string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", hashToken(rawToken)).
		Update("revoked_at", now)
	err := result.Error
	if err != nil {
		return fmt.Errorf("repository.RevokeRefreshToken: %w", err)
	}

	if result.RowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}
	return nil
}

func (r *Repository) RevokeAllRefreshTokens(ctx context.Context, academyID string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("academy_id = ? AND revoked_at IS NULL", academyID).
		Update("revoked_at", now).Error
	if err != nil {
		return fmt.Errorf("repository.RevokeAllRefreshTokens: %w", err)
	}
	return nil
}
