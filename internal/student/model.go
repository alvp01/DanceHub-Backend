package student

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	AcademyID   string    `json:"academy_id" gorm:"column:academy_id;type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"column:name;not null"`
	LastName    string    `json:"last_name" gorm:"column:last_name;not null"`
	Email       string    `json:"email" gorm:"column:email;not null;uniqueIndex"`
	Phone       string    `json:"phone" gorm:"column:phone;not null;uniqueIndex"`
	IdDocument  string    `json:"id_document" gorm:"column:id_document;not null;uniqueIndex"`
	BirthDate   string    `json:"birth_date" gorm:"column:birth_date;not null"`
	Address     string    `json:"address" gorm:"column:address;not null"`
	Allergies   string    `json:"allergies" gorm:"column:allergies"`
	Pathologies string    `json:"pathologies" gorm:"column:pathologies"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type CreateStudentRequest struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IdDocument  string `json:"id_document"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Allergies   string `json:"allergies"`
	Pathologies string `json:"pathologies"`
	AcademyID   string `json:"academy_id"`
}

type CreateStudentResponse struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IdDocument  string `json:"id_document"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Allergies   string `json:"allergies"`
	Pathologies string `json:"pathologies"`
	AcademyID   string `json:"academy_id"`
}

type UpdateStudentRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IdDocument  string `json:"id_document"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Allergies   string `json:"allergies"`
	Pathologies string `json:"pathologies"`
	AcademyID   string `json:"academy_id"`
}

type UpdateStudentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IdDocument  string `json:"id_document"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Allergies   string `json:"allergies"`
	Pathologies string `json:"pathologies"`
	AcademyID   string `json:"academy_id"`
}

type DeleteStudentRequest struct {
	ID string `json:"id"`
}

type DeleteStudentResponse struct {
	Message string `json:"message"`
}

type GetStudentRequest struct {
	ID string `json:"id"`
}

type GetStudentResponse struct {
	ID          string `json:"id"`
	AcademyID   string `json:"academy_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	IdDocument  string `json:"id_document"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
	Allergies   string `json:"allergies"`
	Pathologies string `json:"pathologies"`
}

type ListStudentsRequest struct {
	AcademyID string `json:"academy_id"`
}

type ListStudentsResponse struct {
	Students []Student `json:"students"`
}

type SearchStudentsRequest struct {
	Query string `json:"query"`
}

type SearchStudentsResponse struct {
	Students []Student `json:"students"`
}

type FindByIdDocumentRequest struct {
	IdDocument string `json:"id_document"`
}

func (s *Student) BeforeCreate(_ *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}
