package assignment_repo

import (
	assignment_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/domain"
)

type AssignmentRepoPort interface {
	StoreAssignment(assignment assignment_domain.Assignment) (assignment_domain.Assignment, error)
	UpdateAssignment(assignment assignment_domain.Assignment) (assignment_domain.Assignment, error)
	DeleteAssignment(id string) error
	FindAssignments() ([]assignment_domain.Assignment, error)
	FindAssignmentById(id string) (assignment_domain.Assignment, error)
	FindAssignmentsBySection(sectionId string) ([]assignment_domain.Assignment, error)
	FindAssignmentsByStudent(studentId string) ([]assignment_domain.Assignment, error)
	FindAssignmentsByInstructor(studentId string) ([]assignment_domain.Assignment, error)
}
