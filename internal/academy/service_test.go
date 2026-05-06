package academy

import (
	"context"
	"errors"
	"testing"

	"github.com/alvp01/DanceHub-Backend/internal/validator"
)

func TestServiceRegister_InvalidPassword(t *testing.T) {
	_, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()
	_, err := service.Register(ctx, RegisterRequest{
		Name:         "Test Academy",
		Email:        "invalidpw@example.com",
		PrimaryPhone: "123123123",
		Password:     "short",
	})
	if err == nil {
		t.Fatalf("expected password validation error")
	}
}

func TestServiceLogin_WrongPassword(t *testing.T) {
	_, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()
	_, err := service.Register(ctx, RegisterRequest{
		Name:         "Login Academy",
		Email:        "login@example.com",
		PrimaryPhone: "321321321",
		Password:     "ABcdef123!^x",
	})
	if err != nil {
		t.Fatalf("registering academy: %v", err)
	}
	_, err = service.Login(ctx, LoginRequest{Email: "login@example.com", Password: "wrongpass"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
func TestServiceRegister_FieldUniqueness(t *testing.T) {
	_, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()
	validReq := RegisterRequest{
		Name:         "Unique Academy",
		Email:        "unique@example.com",
		PrimaryPhone: "123456789",
		Password:     "ABcdef123!^x",
	}
	_, err := service.Register(ctx, validReq)
	if err != nil {
		t.Fatalf("unexpected error registering academy: %v", err)
	}

	// Email uniqueness
	_, err = service.Register(ctx, RegisterRequest{
		Name:         "Another Name",
		Email:        validReq.Email,
		PrimaryPhone: "987654321",
		Password:     "ABcdef123!^x",
	})
	if err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}

	// Name uniqueness
	_, err = service.Register(ctx, RegisterRequest{
		Name:         validReq.Name,
		Email:        "other@example.com",
		PrimaryPhone: "111222333",
		Password:     "ABcdef123!^x",
	})
	if err != ErrNameAlreadyExists {
		t.Fatalf("expected ErrNameAlreadyExists, got %v", err)
	}

	// Phone uniqueness
	_, err = service.Register(ctx, RegisterRequest{
		Name:         "Other Academy",
		Email:        "other2@example.com",
		PrimaryPhone: validReq.PrimaryPhone,
		Password:     "ABcdef123!^x",
	})
	if err != ErrPhoneAlreadyExists {
		t.Fatalf("expected ErrPhoneAlreadyExists, got %v", err)
	}
}

func TestServiceRegisterValidationAndLoginFlow(t *testing.T) {
	repo, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()

	_, err := service.Register(ctx, RegisterRequest{
		Name:         "Dance Hub",
		Email:        "academy@example.com",
		PrimaryPhone: "999111222",
		Password:     "short",
	})
	if err == nil {
		t.Fatalf("expected password validation error")
	}

	var pwErr *validator.PasswordValidationError
	if ok := errors.As(err, &pwErr); !ok {
		t.Fatalf("expected PasswordValidationError, got %T", err)
	}

	registerResp, err := service.Register(ctx, RegisterRequest{
		Name:         "Dance Hub",
		Email:        "academy@example.com",
		PrimaryPhone: "999111222",
		Password:     "ABcdef123!^x",
	})
	if err != nil {
		t.Fatalf("registering academy: %v", err)
	}

	if registerResp.ID == "" {
		t.Fatalf("expected non-empty academy id")
	}

	createdAcademy, err := repo.FindByEmail(ctx, "academy@example.com")
	if err != nil {
		t.Fatalf("finding created academy: %v", err)
	}
	if createdAcademy.PasswordHash == "ABcdef123!^x" || createdAcademy.PasswordHash == "" {
		t.Fatalf("expected stored password hash, got invalid value")
	}

	_, err = service.Login(ctx, LoginRequest{Email: "academy@example.com", Password: "wrong"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}

	loginResp, err := service.Login(ctx, LoginRequest{Email: "academy@example.com", Password: "ABcdef123!^x"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginResp.AccessToken == "" || loginResp.RefreshToken == "" {
		t.Fatalf("expected login tokens, got %+v", loginResp)
	}

	refreshResp, err := service.RefreshTokens(ctx, RefreshRequest{RefreshToken: loginResp.RefreshToken})
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	if refreshResp.RefreshToken == loginResp.RefreshToken {
		t.Fatalf("expected new refresh token to differ from previous token")
	}

	_, err = service.RefreshTokens(ctx, RefreshRequest{RefreshToken: loginResp.RefreshToken})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected revoked refresh token to be rejected with invalid credentials, got %v", err)
	}
}

func TestServiceLogoutRevokesRefreshToken(t *testing.T) {
	_, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()

	registerResp, err := service.Register(ctx, RegisterRequest{
		Name:         "Steps Academy",
		Email:        "steps@example.com",
		PrimaryPhone: "999333444",
		Password:     "ABcdef123!^x",
	})
	if err != nil {
		t.Fatalf("registering academy: %v", err)
	}

	loginResp, err := service.Login(ctx, LoginRequest{Email: "steps@example.com", Password: "ABcdef123!^x"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if err := service.Logout(ctx, registerResp.ID, loginResp.RefreshToken); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err = service.RefreshTokens(ctx, RefreshRequest{RefreshToken: loginResp.RefreshToken})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials after logout, got %v", err)
	}
}

func TestServiceLogout_RejectsAnotherAcademyToken(t *testing.T) {
	_, service, _ := setupAcademyTestDeps(t)
	ctx := context.Background()

	first, err := service.Register(ctx, RegisterRequest{
		Name:         "First Academy",
		Email:        "first@example.com",
		PrimaryPhone: "900000001",
		Password:     "ABcdef123!^x",
	})
	if err != nil {
		t.Fatalf("registering first academy: %v", err)
	}

	_, err = service.Register(ctx, RegisterRequest{
		Name:         "Second Academy",
		Email:        "second@example.com",
		PrimaryPhone: "900000002",
		Password:     "ABcdef123!^x",
	})
	if err != nil {
		t.Fatalf("registering second academy: %v", err)
	}

	secondLogin, err := service.Login(ctx, LoginRequest{Email: "second@example.com", Password: "ABcdef123!^x"})
	if err != nil {
		t.Fatalf("logging second academy: %v", err)
	}

	err = service.Logout(ctx, first.ID, secondLogin.RefreshToken)
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials when academy logs out another academy token, got %v", err)
	}

	// Token should still be valid for its owner after rejected logout attempt.
	if _, err := service.RefreshTokens(ctx, RefreshRequest{RefreshToken: secondLogin.RefreshToken}); err != nil {
		t.Fatalf("expected second academy token to remain valid, got %v", err)
	}
}
