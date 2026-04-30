package academy

import "testing"

func TestAcademyBeforeCreateSetsID(t *testing.T) {
	a := &Academy{}

	if err := a.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if a.ID == "" {
		t.Fatalf("expected generated academy id")
	}
}

func TestRefreshTokenBeforeCreateSetsID(t *testing.T) {
	rt := &RefreshToken{}

	if err := rt.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if rt.ID == "" {
		t.Fatalf("expected generated refresh token id")
	}
}
