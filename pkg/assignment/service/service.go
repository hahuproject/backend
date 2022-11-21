package assignment_service

import (
	"log"

	assignment_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/domain"
	assignment_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/repo"
	assignment_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/utils"
)

type AssignmentServiceRepo interface {
	GetAssignments(token string) ([]assignment_domain.Assignment, error)
	AddAssignment(token string, assignment assignment_domain.Assignment) (assignment_domain.Assignment, error)
	UpdateAssignment(token string, assignment assignment_domain.Assignment) (assignment_domain.Assignment, error)
	DeleteAssignment(token string, assignmentId string) error
}

type AssignmentServiceAdapter struct {
	log  *log.Logger
	repo assignment_repo.AssignmentRepoPort
}

func NewAssignmentServiceAdapter(log *log.Logger, repo assignment_repo.AssignmentRepoPort) AssignmentServiceRepo {
	return AssignmentServiceAdapter{log: log, repo: repo}
}

func (service AssignmentServiceAdapter) GetAssignments(token string) ([]assignment_domain.Assignment, error) {

	var _assignments []assignment_domain.Assignment = make([]assignment_domain.Assignment, 0)

	user, err := assignment_utils.CheckAuth(token)

	if err != nil {
		return _assignments, err
	}

	if user.Type == "ADMIN" || user.Type == "REGISTRY_OFFICER" {
		return service.repo.FindAssignments()
	}

	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" {
		return service.repo.FindAssignmentsByInstructor(user.ID)
	}

	if user.Type == "STUDENT" {
		return service.repo.FindAssignmentsByStudent(user.ID)
	}

	return _assignments, assignment_utils.ErrNotAuthorized

}
func (service AssignmentServiceAdapter) AddAssignment(token string, assignment assignment_domain.Assignment) (assignment_domain.Assignment, error) {

	var _addedAssignment assignment_domain.Assignment

	_, err := assignment_utils.CheckAuth(token)

	if err != nil {
		return _addedAssignment, err
	}

	return service.repo.StoreAssignment(assignment)

}
func (service AssignmentServiceAdapter) UpdateAssignment(token string, assignment assignment_domain.Assignment) (assignment_domain.Assignment, error) {

	var _updatedAssignment assignment_domain.Assignment

	_, err := assignment_utils.CheckAuth(token)

	if err != nil {
		return _updatedAssignment, err
	}

	return service.repo.UpdateAssignment(assignment)

}
func (service AssignmentServiceAdapter) DeleteAssignment(token string, assignmentId string) error {

	_, err := assignment_utils.CheckAuth(token)

	if err != nil {
		return err
	}

	return service.repo.DeleteAssignment(assignmentId)

}
