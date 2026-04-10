// internal/academy/service.go
package academy

import (
	"context"
	"fmt"
	"os"
	"strconv"

	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
	"github.com/alvp01/DanceHub-Backend/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       *Repository
	jwtManager *jwtpkg.Manager
}

func NewService(repo *Repository, jwtManager *jwtpkg.Manager) *Service {
	return &Service{repo: repo, jwtManager: jwtManager}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// 1. Validar password
	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// 2. Validar campos requeridos básicos
	if req.Name == "" || req.Email == "" || req.PrimaryPhone == "" {
		return nil, fmt.Errorf("nombre, email y teléfono primario son obligatorios")
	}

	// 3. Hashear el password con bcrypt
	cost, _ := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		return nil, fmt.Errorf("service.Register: error hasheando password: %w", err)
	}

	// 4. Construir entidad
	academy := &Academy{
		Name:           req.Name,
		Email:          req.Email,
		PrimaryPhone:   req.PrimaryPhone,
		SecondaryPhone: req.SecondaryPhone,
		PasswordHash:   string(hash),
	}

	// 5. Persistir
	if err := s.repo.Create(ctx, academy); err != nil {
		return nil, err
	}

	// 6. Responder sin datos sensibles
	return &RegisterResponse{
		ID:             academy.ID,
		Name:           academy.Name,
		Email:          academy.Email,
		PrimaryPhone:   academy.PrimaryPhone,
		SecondaryPhone: academy.SecondaryPhone,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email y password son requeridos")
	}

	// 1. Buscar academia por email
	academy, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// 2. Verificar password contra el hash
	if err := bcrypt.CompareHashAndPassword(
		[]byte(academy.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Generar access token
	accessToken, err := s.jwtManager.GenerateAccessToken(academy.ID, academy.Name)
	if err != nil {
		return nil, fmt.Errorf("service.Login: error generando access token: %w", err)
	}

	// 4. Generar refresh token
	rawRefreshToken, expiresAt, err := s.jwtManager.GenerateRefreshToken(academy.ID, academy.Name)
	if err != nil {
		return nil, fmt.Errorf("service.Login: error generando refresh token: %w", err)
	}

	// 5. Persistir el refresh token (hasheado)
	if err := s.repo.SaveRefreshToken(ctx, academy.ID, rawRefreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.AccessDuration().Seconds()),
	}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, req RefreshRequest) (*LoginResponse, error) {
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token es requerido")
	}

	// 1. Validar firma y expiración del JWT
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verificar que existe en DB y no está revocado
	_, err = s.repo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Rotación: revocar el token usado (refresh token rotation)
	if err := s.repo.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	// 4. Generar nuevos tokens
	newAccessToken, err := s.jwtManager.GenerateAccessToken(claims.AcademyID, claims.AcademyName)
	if err != nil {
		return nil, fmt.Errorf("service.RefreshTokens: %w", err)
	}

	newRawRefresh, expiresAt, err := s.jwtManager.GenerateRefreshToken(claims.AcademyID, claims.AcademyName)
	if err != nil {
		return nil, fmt.Errorf("service.RefreshTokens: %w", err)
	}

	if err := s.repo.SaveRefreshToken(ctx, claims.AcademyID, newRawRefresh, expiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.AccessDuration().Seconds()),
	}, nil
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, rawRefreshToken)
}

func (s *Service) LogoutAll(ctx context.Context, academyID string) error {
	return s.repo.RevokeAllRefreshTokens(ctx, academyID)
}
