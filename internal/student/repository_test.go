package student

import (
	"context"
	"testing"
)

func TestRepositoryCreateStudent(t *testing.T) {
	repo, _, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	created := &Student{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana@example.com",
		Phone:       "300111222",
		IdDocument:  "DOC-100",
		BirthDate:   "2000-01-01",
		Address:     "Street 1",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("creating student: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected generated student id")
	}
}

func TestRepositoryReadStudentByIDAndDocument(t *testing.T) {
	repo, _, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	created := &Student{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.read@example.com",
		Phone:       "300111232",
		IdDocument:  "DOC-110",
		BirthDate:   "2000-01-01",
		Address:     "Street 1",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("creating student: %v", err)
	}

	foundByDoc, err := repo.FindByIdDocument(ctx, "DOC-110")
	if err != nil {
		t.Fatalf("finding by document: %v", err)
	}
	if foundByDoc.ID != created.ID {
		t.Fatalf("expected same student id, got %q and %q", foundByDoc.ID, created.ID)
	}

	foundByID, err := repo.FindById(ctx, created.ID)
	if err != nil {
		t.Fatalf("finding by id: %v", err)
	}
	if foundByID.Email != "ana.read@example.com" {
		t.Fatalf("unexpected email: %s", foundByID.Email)
	}
}

func TestRepositoryUpdateStudent(t *testing.T) {
	repo, _, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	created := &Student{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.update@example.com",
		Phone:       "300111242",
		IdDocument:  "DOC-120",
		BirthDate:   "2000-01-01",
		Address:     "Street 1",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("creating student: %v", err)
	}

	updatePayload := &Student{
		ID:          created.ID,
		Name:        "Ana Maria",
		LastName:    "Lopez",
		Email:       "ana.maria@example.com",
		Phone:       "300111223",
		IdDocument:  "DOC-101",
		BirthDate:   "2000-01-01",
		Address:     "Street 2",
		Allergies:   "peanuts",
		Pathologies: "none",
	}
	if err := repo.Update(ctx, updatePayload); err != nil {
		t.Fatalf("updating student: %v", err)
	}

	updated, err := repo.FindById(ctx, created.ID)
	if err != nil {
		t.Fatalf("finding updated student: %v", err)
	}
	if updated.Name != "Ana Maria" || updated.IdDocument != "DOC-101" {
		t.Fatalf("student was not updated correctly: %+v", updated)
	}
}

func TestRepositoryDeleteStudent(t *testing.T) {
	repo, _, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	created := &Student{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.delete@example.com",
		Phone:       "300111252",
		IdDocument:  "DOC-130",
		BirthDate:   "2000-01-01",
		Address:     "Street 1",
		Allergies:   "none",
		Pathologies: "none",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("creating student: %v", err)
	}

	all, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("listing students: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 student, got %d", len(all))
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("deleting student: %v", err)
	}

	_, err = repo.FindById(ctx, created.ID)
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound after delete, got %v", err)
	}
}

func TestRepositoryNotFoundCases(t *testing.T) {
	repo, _, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := repo.FindById(ctx, "missing-id")
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound, got %v", err)
	}

	_, err = repo.FindByIdDocument(ctx, "missing-doc")
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound, got %v", err)
	}

	err = repo.Update(ctx, &Student{ID: "missing-id", Name: "X"})
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound on update, got %v", err)
	}

	err = repo.Delete(ctx, "missing-id")
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound on delete, got %v", err)
	}
}
