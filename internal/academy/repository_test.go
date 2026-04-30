package academy

import (
	"context"
	"testing"
	"time"
)

func TestRepositoryFindByEmailNotFound(t *testing.T) {
	repo, _, _ := setupAcademyTestDeps(t)

	_, err := repo.FindByEmail(context.Background(), "missing@example.com")
	if err != ErrAcademyNotFound {
		t.Fatalf("expected ErrAcademyNotFound, got %v", err)
	}
}

func TestRepositoryFindRefreshTokenRevokedAndExpired(t *testing.T) {
	repo, _, _ := setupAcademyTestDeps(t)
	ctx := context.Background()

	if err := repo.SaveRefreshToken(ctx, "academy-1", "revoked-token", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("saving revoked token fixture: %v", err)
	}
	if err := repo.RevokeRefreshToken(ctx, "revoked-token"); err != nil {
		t.Fatalf("revoking token fixture: %v", err)
	}

	_, err := repo.FindRefreshToken(ctx, "revoked-token")
	if err != ErrRefreshTokenRevoked {
		t.Fatalf("expected ErrRefreshTokenRevoked, got %v", err)
	}

	if err := repo.SaveRefreshToken(ctx, "academy-1", "expired-token", time.Now().Add(-time.Minute)); err != nil {
		t.Fatalf("saving expired token fixture: %v", err)
	}

	_, err = repo.FindRefreshToken(ctx, "expired-token")
	if err != ErrRefreshTokenExpired {
		t.Fatalf("expected ErrRefreshTokenExpired, got %v", err)
	}
}
