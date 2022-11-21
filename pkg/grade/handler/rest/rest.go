package grade_rest_handler

import (
	"encoding/json"
	"log"
	"net/http"

	grade_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/domain"
	grade_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/handler"
	grade_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/service"
	grade_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/utils"
)

type GradeRestHandlerAdapter struct {
	log     *log.Logger
	service grade_service.GradeServicePort
}

func NewGradeRestHandlerAdapter(log *log.Logger, service grade_service.GradeServicePort) grade_handler.GradeHandlerPort {
	return GradeRestHandlerAdapter{log: log, service: service}
}

func (handler GradeRestHandlerAdapter) GetAddGrade(w http.ResponseWriter, r *http.Request) {
	// print("grade handler 0")
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	// print("grade handler 1")
	if err != nil {
		// print("grade handler 2")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// print("grade handler 3")
	var grade grade_domain.Grade

	// print("grade handler 4")
	decoder := json.NewDecoder(r.Body)

	// print("grade handler 5")
	err = decoder.Decode(&grade)

	// print("grade handler 6")
	if err != nil {
		// print("grade handler 7")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// print("grade handler 8")
	defer r.Body.Close()
	// print("grade handler 9")
	// print(grade)
	handler.log.Println(grade)

	addedGrade, err := handler.service.AddGrade(token, grade)
	// print("grade handler 10")
	if err != nil {
		// print("grade handler 11")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// print("grade handler 12")
	grade_utils.SendMarshaledResponse(w, addedGrade)
	// print("grade handler 13")
}
func (handler GradeRestHandlerAdapter) GetGrades(w http.ResponseWriter, r *http.Request) {

	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	grades, err := handler.service.GetGrades(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	grade_utils.SendMarshaledResponse(w, grades)
}

func (handler GradeRestHandlerAdapter) GetRequestGradeReview(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	courseId := r.URL.Query().Get("course-id")
	userId := r.URL.Query().Get("user-id")

	if courseId == "" || userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameters not specified"))
		return

	}

	err = handler.service.RequestGradeReview(token, courseId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully submitted grade review request"))
}

func (handler GradeRestHandlerAdapter) GetApproveGradeReview(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	courseId := r.URL.Query().Get("course-id")
	userId := r.URL.Query().Get("user-id")

	if courseId == "" || userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameters not specified"))
		return

	}

	err = handler.service.ApproveGradeReview(token, userId, courseId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully approved grade review request"))
}

func (handler GradeRestHandlerAdapter) GetRejectGradeReview(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	courseId := r.URL.Query().Get("course-id")
	userId := r.URL.Query().Get("user-id")

	if courseId == "" || userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameters not specified"))
		return

	}

	err = handler.service.RejectGradeReview(token, userId, courseId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully rejected grade review request"))
}

func (handler GradeRestHandlerAdapter) GetSubmitGrade(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	courseId := r.URL.Query().Get("course-id")
	userId := r.URL.Query().Get("user-id")

	if courseId == "" || userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameters not specified"))
		return

	}

	err = handler.service.SubmitGrade(token, userId, courseId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully aubmitted grade"))
}

//Grade Labels
func (handler GradeRestHandlerAdapter) GetAddGradeLabel(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not Authorized"))
		return
	}

	var gradeLabel grade_domain.GradeLabel

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&gradeLabel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedGradeLabel, err := handler.service.AddGradeLabel(token, gradeLabel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	grade_utils.SendMarshaledResponse(w, addedGradeLabel)

}
func (handler GradeRestHandlerAdapter) GetGradeLabels(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not Authorized"))
		return
	}

	gradeLabels, err := handler.service.GetGradeLabels(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	grade_utils.SendMarshaledResponse(w, gradeLabels)

}
func (handler GradeRestHandlerAdapter) GetRemoveGradeLabels(w http.ResponseWriter, r *http.Request) {
	token, err := grade_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not Authorized"))
		return
	}

	gradeLabelId := r.URL.Query().Get("id")

	err = handler.service.RemoveGradeLabel(token, gradeLabelId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted grade label"))
}
