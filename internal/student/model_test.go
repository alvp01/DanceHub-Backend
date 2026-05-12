package student

import "testing"

func TestStudentBeforeCreateSetsID(t *testing.T) {
	s := &Student{}

	if err := s.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if s.ID == "" {
		t.Fatalf("expected generated student id")
	}
}
