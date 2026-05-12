package student

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type StudenError string

func (e StudenError) Error() string {
	return string(e)
}

const (
	ErrStudentNotFound    StudenError = "estudiante no encontrado"
	ErrEmailAlreadyExists StudenError = "el email ya está registrado"
	ErrPhoneAlreadyExists StudenError = "el teléfono ya está registrado"
	ErrIdDocAlreadyExists StudenError = "el documento de identidad ya está registrado"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, s *Student) error {
	err := r.db.WithContext(ctx).Create(s).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "students_email_key"):
			return ErrEmailAlreadyExists
		case strings.Contains(err.Error(), "students_phone_key"):
			return ErrPhoneAlreadyExists
		case strings.Contains(err.Error(), "students_id_document_key"):
			return ErrIdDocAlreadyExists
		}
		return fmt.Errorf("repository.Create: %w", err)
	}

	return nil
}

func (r *Repository) FindById(ctx context.Context, id string) (*Student, error) {
	s := &Student{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("repository.FindById: %w", err)
	}

	return s, nil
}

func (r *Repository) Update(ctx context.Context, s *Student) error {
	result := r.db.WithContext(ctx).
		Model(&Student{}).
		Where("id = ?", s.ID).
		Updates(map[string]any{
			"name":        s.Name,
			"last_name":   s.LastName,
			"email":       s.Email,
			"phone":       s.Phone,
			"id_document": s.IdDocument,
			"birth_date":  s.BirthDate,
			"address":     s.Address,
			"allergies":   s.Allergies,
			"pathologies": s.Pathologies,
		})

	if result.Error != nil {
		return fmt.Errorf("repository.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Student{})
	if result.Error != nil {
		return fmt.Errorf("repository.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*Student, error) {
	var students []*Student
	err := r.db.WithContext(ctx).Find(&students).Error
	if err != nil {
		return nil, fmt.Errorf("repository.FindAll: %w", err)
	}

	return students, nil
}

func (r *Repository) FindByIdDocument(ctx context.Context, idDocument string) (*Student, error) {
	s := &Student{}
	err := r.db.WithContext(ctx).Where("id_document = ?", idDocument).First(s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("repository.FindByIdDocument: %w", err)
	}

	return s, nil
}
