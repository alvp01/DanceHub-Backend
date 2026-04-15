// internal/academy/repository.go
package academy

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *Academy) error {
	query := `
        INSERT INTO academies (name, email, primary_phone, secondary_phone, password_hash)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx, query,
		a.Name,
		a.Email,
		a.PrimaryPhone,
		a.SecondaryPhone,
		a.PasswordHash,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		// Detectar violaciones de unicidad de Postgres
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "academies_email_key":
				return ErrEmailAlreadyExists
			case "academies_name_key":
				return ErrNameAlreadyExists
			case "academies_primary_phone_key":
				return ErrPhoneAlreadyExists
			}
		}
		return fmt.Errorf("repository.Create: %w", err)
	}

	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*Academy, error) {
	query := `
        SELECT id, name, email, primary_phone, secondary_phone,
               password_hash, created_at, updated_at
        FROM academies
        WHERE email = $1
    `

	a := &Academy{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&a.ID, &a.Name, &a.Email,
		&a.PrimaryPhone, &a.SecondaryPhone,
		&a.PasswordHash, &a.CreatedAt, &a.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
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
	query := `
        INSERT INTO refresh_tokens (academy_id, token_hash, expires_at)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.ExecContext(ctx, query, academyID, hashToken(rawToken), expiresAt)
	if err != nil {
		return fmt.Errorf("repository.SaveRefreshToken: %w", err)
	}
	return nil
}

func (r *Repository) FindRefreshToken(ctx context.Context, rawToken string) (*RefreshToken, error) {
	query := `
        SELECT id, academy_id, token_hash, expires_at, created_at, revoked_at
        FROM refresh_tokens
        WHERE token_hash = $1
    `

	rt := &RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, hashToken(rawToken)).Scan(
		&rt.ID, &rt.AcademyID, &rt.TokenHash,
		&rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
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

// RevokeRefreshToken invalida un token específico (logout)
func (r *Repository) RevokeRefreshToken(ctx context.Context, rawToken string) error {
	query := `
        UPDATE refresh_tokens
        SET revoked_at = NOW()
        WHERE token_hash = $1 AND revoked_at IS NULL
    `
	result, err := r.db.ExecContext(ctx, query, hashToken(rawToken))
	if err != nil {
		return fmt.Errorf("repository.RevokeRefreshToken: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrRefreshTokenNotFound
	}
	return nil
}

// RevokeAllRefreshTokens invalida todos los tokens de una academia (logout de todos los dispositivos)
func (r *Repository) RevokeAllRefreshTokens(ctx context.Context, academyID string) error {
	query := `
        UPDATE refresh_tokens
        SET revoked_at = NOW()
        WHERE academy_id = $1 AND revoked_at IS NULL
    `
	_, err := r.db.ExecContext(ctx, query, academyID)
	if err != nil {
		return fmt.Errorf("repository.RevokeAllRefreshTokens: %w", err)
	}
	return nil
}
