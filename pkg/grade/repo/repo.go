package grade_repo

import (
	grade_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/domain"
)

type GradeRepoPort interface {
	StoreGrade(grade grade_domain.Grade) (grade_domain.Grade, error)
	FindGrades() ([]grade_domain.Grade, error)
	FindGradesByInstructor(instructorId string) ([]grade_domain.Grade, error)
	FindGradesByUser(userId string) ([]grade_domain.Grade, error)
	RequestGradeReview(userId, courseId string) error
	ApproveGradeReview(userId, courseId string) error
	RejectGradeReview(userId, courseId string) error
	SubmitGrade(userId, courseId string) error

	//Grade Label
	StoreGradeLabel(gradeLabel grade_domain.GradeLabel) (grade_domain.GradeLabel, error)
	FindGradeLabels() ([]grade_domain.GradeLabel, error)
	RemoveGradeLabels(gradeLabelId string) error
}
