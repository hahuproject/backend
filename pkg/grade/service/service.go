package grade_service

import (
	"log"

	grade_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/domain"
	grade_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/repo"
	grade_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/utils"
)

type GradeServicePort interface {
	AddGrade(token string, grade grade_domain.Grade) (grade_domain.Grade, error)
	GetGrades(token string) ([]grade_domain.Grade, error)
	RequestGradeReview(tokens, courseId, userId string) error
	ApproveGradeReview(token, userId, courseId string) error
	RejectGradeReview(token, userId, courseId string) error
	SubmitGrade(token, userId, courseId string) error

	//Grade Labels
	AddGradeLabel(token string, gradeLabel grade_domain.GradeLabel) (grade_domain.GradeLabel, error)
	GetGradeLabels(token string) ([]grade_domain.GradeLabel, error)
	RemoveGradeLabel(token, gradeLabelId string) error
}

type GradeServiceAdapter struct {
	log  *log.Logger
	repo grade_repo.GradeRepoPort
}

func NewGradeServiceAdapter(log *log.Logger, repo grade_repo.GradeRepoPort) GradeServicePort {
	return GradeServiceAdapter{log: log, repo: repo}
}

func (service GradeServiceAdapter) AddGrade(token string, grade grade_domain.Grade) (grade_domain.Grade, error) {

	var addedGrade grade_domain.Grade

	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return addedGrade, err
	}

	if user.Type != "INSTRUCTOR" && user.Type != "SUPERVISOR" {
		return addedGrade, grade_utils.ErrNotAuthorized
	}

	return service.repo.StoreGrade(grade)
}
func (service GradeServiceAdapter) GetGrades(token string) ([]grade_domain.Grade, error) {

	var grades []grade_domain.Grade = make([]grade_domain.Grade, 0)

	// print("grade serv 0")

	user, err := grade_utils.CheckAuth(token)
	// print("grade serv 1")
	if err != nil {
		// print("grade serv 2")
		return grades, err
	}

	// service.log.Println("grade serv 3")
	// service.log.Println(user.Type)
	// service.log.Println(user.ID)
	if user.Type == "STUDENT" {
		// print("grade serv 4")
		return service.repo.FindGradesByUser(user.ID)
	}

	// print("grade serv 5")
	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" {
		// print("grade serv 6")
		return service.repo.FindGradesByInstructor(user.ID)
	}

	if user.Type == "SUPER_ADMIN" || user.Type == "ADMIN" || user.Type == "REGISTRY_OFFICER" || user.Type == "DEPARTMENT_HEAD" {
		// print("grade serv 6")
		return service.repo.FindGrades()
	}

	// print("grade serv 7")
	return grades, grade_utils.ErrNotAuthorized
}

func (service GradeServiceAdapter) RequestGradeReview(token, courseId, userId string) error {
	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "STUDENT" && user.Type != "INSTRUCTOR" && user.Type != "SUPERVISOR" {
		return grade_utils.ErrNotAuthorized
	}

	return service.repo.RequestGradeReview(userId, courseId)
}
func (service GradeServiceAdapter) ApproveGradeReview(token, userId, courseId string) error {
	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return grade_utils.ErrNotAuthorized
	}

	return service.repo.ApproveGradeReview(userId, courseId)
}
func (service GradeServiceAdapter) RejectGradeReview(token, userId, courseId string) error {
	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return grade_utils.ErrNotAuthorized
	}

	return service.repo.RejectGradeReview(userId, courseId)
}
func (service GradeServiceAdapter) SubmitGrade(token, userId, courseId string) error {
	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "INSTRUCTOR" && user.Type != "SUPERVISOR" {
		return grade_utils.ErrNotAuthorized
	}

	return service.repo.SubmitGrade(userId, courseId)
}

//Grade Labels
func (service GradeServiceAdapter) AddGradeLabel(token string, gradeLabel grade_domain.GradeLabel) (grade_domain.GradeLabel, error) {

	var addedGradeLabel grade_domain.GradeLabel

	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return addedGradeLabel, err
	}

	if user.Type != "ADMIN" && user.Type != "SUPER_ADMIN" {
		return addedGradeLabel, grade_utils.ErrNotAuthorized
	}

	return service.repo.StoreGradeLabel(gradeLabel)

}
func (service GradeServiceAdapter) GetGradeLabels(token string) ([]grade_domain.GradeLabel, error) {

	var gradeLabels []grade_domain.GradeLabel = make([]grade_domain.GradeLabel, 0)

	_, err := grade_utils.CheckAuth(token)
	if err != nil {
		return gradeLabels, err
	}

	// if user.Type != "ADMIN" && user.Type != "SUPER_ADMIN" {
	// 	return gradeLabels, grade_utils.ErrNotAuthorized
	// }

	return service.repo.FindGradeLabels()

}
func (service GradeServiceAdapter) RemoveGradeLabel(token, gradeLabelId string) error {
	user, err := grade_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "ADMIN" && user.Type != "SUPER_ADMIN" {
		return grade_utils.ErrNotAuthorized
	}

	return service.repo.RemoveGradeLabels(gradeLabelId)
}
