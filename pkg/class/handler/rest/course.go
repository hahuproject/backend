package class_rest_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (courseHandler RestClassHandlerAdapter) GetAddCourse(w http.ResponseWriter, r *http.Request) {
	var course class_domain.Course

	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&course)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured while decoding body"))
		return
	}
	defer r.Body.Close()

	addedCourse, err := courseHandler.classService.AddCourse(strings.Split(token, " ")[1], course)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledAddedCourse, _ := json.Marshal(addedCourse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedCourse)
}
func (courseHandler RestClassHandlerAdapter) GetCourses(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	courses, err := courseHandler.classService.GetCourses(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledCourses, _ := json.Marshal(courses)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledCourses)

}
func (courseHandler RestClassHandlerAdapter) GetCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	course, err := courseHandler.classService.GetCourse(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledCourse, _ := json.Marshal(course)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledCourse)
}
func (courseHandler RestClassHandlerAdapter) GetUpdateCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	var course class_domain.Course

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&course)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured while decoding body"))
		return
	}
	defer r.Body.Close()

	updatedCourse, err := courseHandler.classService.UpdateCourse(token, course)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUpdatedCourse, _ := json.Marshal(updatedCourse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUpdatedCourse)
}
func (courseHandler RestClassHandlerAdapter) GetDeleteCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	err := courseHandler.classService.DeleteCourse(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted course"))
}
func (courseHandler RestClassHandlerAdapter) GetAddUserToCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	type CourseUserIds struct {
		CourseId string `json:"courseId"`
		UserId   string `json:"userId"`
	}

	var courseUserId CourseUserIds

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&courseUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding body"))
		return
	}

	defer r.Body.Close()

	updatedCourse, err := courseHandler.classService.AddUserToCourse(token, courseUserId.CourseId, courseUserId.UserId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, updatedCourse)

}
func (courseHandler RestClassHandlerAdapter) GetRemoveUserFromCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	type CourseUserIds struct {
		CourseId string `json:"courseId"`
		UserId   string `json:"userId"`
	}

	var courseUserId CourseUserIds

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&courseUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding body"))
		return
	}

	defer r.Body.Close()

	updatedCourse, err := courseHandler.classService.RemoveUserFromCourse(token, courseUserId.CourseId, courseUserId.UserId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, updatedCourse)
}
