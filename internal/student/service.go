package student

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateStudentRequest) (*CreateStudentResponse, error) {

	if req.Name == "" || req.LastName == "" || req.Email == "" {
		return nil, fmt.Errorf("nombre, apellido y email son obligatorios")
	}

	student := &Student{
		AcademyID:   req.AcademyID,
		Name:        req.Name,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		IdDocument:  req.IdDocument,
		BirthDate:   req.BirthDate,
		Address:     req.Address,
		Allergies:   req.Allergies,
		Pathologies: req.Pathologies,
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return nil, err
	}

	return &CreateStudentResponse{
		Name:        student.Name,
		LastName:    student.LastName,
		Email:       student.Email,
		Phone:       student.Phone,
		IdDocument:  student.IdDocument,
		BirthDate:   student.BirthDate,
		Address:     student.Address,
		Allergies:   student.Allergies,
		Pathologies: student.Pathologies,
	}, nil
}

func (s *Service) GetById(ctx context.Context, id string) (*Student, error) {
	if id == "" {
		return nil, fmt.Errorf("el ID del estudiante es obligatorio")
	}

	student, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Service) Update(ctx context.Context, req UpdateStudentRequest) (*UpdateStudentResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("el ID del estudiante es obligatorio")
	}

	student := &Student{
		ID:          req.ID,
		Name:        req.Name,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		IdDocument:  req.IdDocument,
		BirthDate:   req.BirthDate,
		Address:     req.Address,
		Allergies:   req.Allergies,
		Pathologies: req.Pathologies,
	}

	if err := s.repo.Update(ctx, student); err != nil {
		return nil, err
	}

	return &UpdateStudentResponse{
		ID:          student.ID,
		Name:        student.Name,
		LastName:    student.LastName,
		Email:       student.Email,
		Phone:       student.Phone,
		IdDocument:  student.IdDocument,
		BirthDate:   student.BirthDate,
		Address:     student.Address,
		Allergies:   student.Allergies,
		Pathologies: student.Pathologies,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("el ID del estudiante es obligatorio")
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]*Student, error) {
	students, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Service) FindByIdDocument(ctx context.Context, idDocument string) (*Student, error) {
	if idDocument == "" {
		return nil, fmt.Errorf("el documento de identidad del estudiante es obligatorio")
	}

	student, err := s.repo.FindByIdDocument(ctx, idDocument)
	if err != nil {
		return nil, err
	}

	return student, nil
}
