package grade_handler

import "net/http"

type GradeHandlerPort interface {
	GetAddGrade(w http.ResponseWriter, r *http.Request)
	GetGrades(w http.ResponseWriter, r *http.Request)
	GetRequestGradeReview(w http.ResponseWriter, r *http.Request)
	GetApproveGradeReview(w http.ResponseWriter, r *http.Request)
	GetRejectGradeReview(w http.ResponseWriter, r *http.Request)
	GetSubmitGrade(w http.ResponseWriter, r *http.Request)

	//Grade Labels
	GetAddGradeLabel(w http.ResponseWriter, r *http.Request)
	GetGradeLabels(w http.ResponseWriter, r *http.Request)
	GetRemoveGradeLabels(w http.ResponseWriter, r *http.Request)
}
