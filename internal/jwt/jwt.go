// internal/jwt/jwt.go
package jwt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims es el payload del JWT
type Claims struct {
	AcademyID   string    `json:"academy_id"`
	AcademyName string    `json:"academy_name"`
	TokenType   TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret    []byte
	refreshSecret   []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewManager() (*Manager, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	if accessSecret == "" || refreshSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET y JWT_REFRESH_SECRET son requeridos")
	}
	if accessSecret == refreshSecret {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET y JWT_REFRESH_SECRET deben ser diferentes")
	}

	accessMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_DURATION_MINUTES"))
	if accessMinutes == 0 {
		accessMinutes = 15
	}

	refreshDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_DURATION_DAYS"))
	if refreshDays == 0 {
		refreshDays = 7
	}

	return &Manager{
		accessSecret:    []byte(accessSecret),
		refreshSecret:   []byte(refreshSecret),
		accessDuration:  time.Duration(accessMinutes) * time.Minute,
		refreshDuration: time.Duration(refreshDays) * 24 * time.Hour,
	}, nil
}

// GenerateAccessToken crea un access token firmado
func (m *Manager) GenerateAccessToken(academyID, academyName string) (string, error) {
	return m.generate(academyID, academyName, AccessToken, m.accessSecret, m.accessDuration)
}

// GenerateRefreshToken crea un refresh token firmado
func (m *Manager) GenerateRefreshToken(academyID, academyName string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.refreshDuration)
	token, err := m.generate(academyID, academyName, RefreshToken, m.refreshSecret, m.refreshDuration)
	return token, expiresAt, err
}

func (m *Manager) generate(
	academyID, academyName string,
	tokenType TokenType,
	secret []byte,
	duration time.Duration,
) (string, error) {
	now := time.Now()

	claims := &Claims{
		AcademyID:   academyID,
		AcademyName: academyName,
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			Issuer:    "academias-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("jwt.generate: %w", err)
	}

	return signed, nil
}

// ValidateAccessToken valida y retorna los claims de un access token
func (m *Manager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return m.validate(tokenStr, m.accessSecret, AccessToken)
}

// ValidateRefreshToken valida y retorna los claims de un refresh token
func (m *Manager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return m.validate(tokenStr, m.refreshSecret, RefreshToken)
}

func (m *Manager) validate(tokenStr string, secret []byte, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		// Verificar que el algoritmo sea el esperado
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma inesperado: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verificar que el tipo de token sea el correcto
	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// AccessDuration expone la duración para usarla en respuestas
func (m *Manager) AccessDuration() time.Duration {
	return m.accessDuration
}
