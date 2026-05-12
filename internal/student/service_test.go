package student

import (
	"context"
	"testing"
)

func TestServiceCreateValidation(t *testing.T) {
	_, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := service.Create(ctx, CreateStudentRequest{
		LastName: "Lopez",
		Email:    "ana@example.com",
	})
	if err == nil {
		t.Fatalf("expected validation error for missing required fields")
	}
}

func TestServiceCreateStudentSuccess(t *testing.T) {
	_, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	created, err := service.Create(ctx, CreateStudentRequest{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.create@example.com",
		Phone:       "300111222",
		IdDocument:  "DOC-200",
		BirthDate:   "2000-01-01",
		Address:     "Street 10",
		Allergies:   "none",
		Pathologies: "none",
	})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}
	if created.Name != "Ana" {
		t.Fatalf("unexpected create response: %+v", created)
	}
}

func TestServiceFindByIdDocument(t *testing.T) {
	_, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := service.Create(ctx, CreateStudentRequest{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.find@example.com",
		Phone:       "300111223",
		IdDocument:  "DOC-201",
		BirthDate:   "2000-01-01",
		Address:     "Street 10",
		Allergies:   "none",
		Pathologies: "none",
	})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}

	byDocument, err := service.FindByIdDocument(ctx, "DOC-201")
	if err != nil {
		t.Fatalf("finding student by document: %v", err)
	}
	if byDocument.Email != "ana.find@example.com" {
		t.Fatalf("unexpected student email: %s", byDocument.Email)
	}
}

func TestServiceFindAll(t *testing.T) {
	_, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := service.Create(ctx, CreateStudentRequest{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana.list@example.com",
		Phone:       "300111224",
		IdDocument:  "DOC-202",
		BirthDate:   "2000-01-01",
		Address:     "Street 10",
		Allergies:   "none",
		Pathologies: "none",
	})
	if err != nil {
		t.Fatalf("creating student: %v", err)
	}

	all, err := service.FindAll(ctx)
	if err != nil {
		t.Fatalf("listing students: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 student, got %d", len(all))
	}
}

func TestServiceUpdateStudent(t *testing.T) {
	repo, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := service.Create(ctx, CreateStudentRequest{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana2@example.com",
		Phone:       "300111224",
		IdDocument:  "DOC-201",
		BirthDate:   "2000-02-02",
		Address:     "Street 20",
		Allergies:   "none",
		Pathologies: "none",
	})
	if err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	stored, err := repo.FindByIdDocument(ctx, "DOC-201")
	if err != nil {
		t.Fatalf("loading fixture student: %v", err)
	}

	updated, err := service.Update(ctx, UpdateStudentRequest{
		ID:          stored.ID,
		Name:        "Ana Maria",
		LastName:    stored.LastName,
		Email:       "ana2-updated@example.com",
		Phone:       "300111225",
		IdDocument:  "DOC-202",
		BirthDate:   stored.BirthDate,
		Address:     "Street 22",
		Allergies:   "dust",
		Pathologies: "none",
	})
	if err != nil {
		t.Fatalf("updating student: %v", err)
	}
	if updated.Name != "Ana Maria" || updated.IdDocument != "DOC-202" {
		t.Fatalf("unexpected update response: %+v", updated)
	}
}

func TestServiceDeleteStudent(t *testing.T) {
	repo, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	fixture := &Student{
		AcademyID:   "academy-1",
		Name:        "Ana",
		LastName:    "Lopez",
		Email:       "ana3@example.com",
		Phone:       "300111226",
		IdDocument:  "DOC-203",
		BirthDate:   "2000-02-02",
		Address:     "Street 20",
		Allergies:   "none",
		Pathologies: "none",
	}
	err := repo.Create(ctx, fixture)
	if err != nil {
		t.Fatalf("creating fixture student: %v", err)
	}

	err = service.Delete(ctx, fixture.ID)
	if err != nil {
		t.Fatalf("deleting student: %v", err)
	}

	err = service.Delete(ctx, fixture.ID)
	if err != ErrStudentNotFound {
		t.Fatalf("expected ErrStudentNotFound on second delete, got %v", err)
	}
}

func TestServiceFindByIdDocumentValidation(t *testing.T) {
	_, service, _ := setupStudentTestDeps(t)
	ctx := context.Background()

	_, err := service.FindByIdDocument(ctx, "")
	if err == nil {
		t.Fatalf("expected validation error when document is empty")
	}
}
