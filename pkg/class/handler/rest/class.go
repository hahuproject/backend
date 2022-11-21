package class_rest_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (classHandler RestClassHandlerAdapter) GetAddClass(w http.ResponseWriter, r *http.Request) {

	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var class class_domain.Class

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedClass, err := classHandler.classService.AddClass(token, class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, addedClass)
}

func (classHandler RestClassHandlerAdapter) GetClasses(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	classes, err := classHandler.classService.GetClasses(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledClasses, err := json.Marshal(classes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledClasses)
}
func (classHandler RestClassHandlerAdapter) GetClass(w http.ResponseWriter, r *http.Request) {

	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(r.URL.Query()["id"]) < 1 || r.URL.Query()["id"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid class id"))
		return
	}

	class, err := classHandler.classService.GetClass(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledClass, err := json.Marshal(class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledClass)
}
func (classHandler RestClassHandlerAdapter) GetUpdateClass(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	var class class_domain.Class

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured while decoding body"))
		return
	}
	defer r.Body.Close()

	updatedClass, err := classHandler.classService.UpdateClass(token, class)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUpdatedCourse, _ := json.Marshal(updatedClass)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUpdatedCourse)
}
func (classHandler RestClassHandlerAdapter) GetDeleteClass(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	err := classHandler.classService.DeleteClass(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted class"))
}

func (classHandler RestClassHandlerAdapter) GetAddClassCourse(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type ClassCourseIds struct {
		ClassId  string `json:"classId"`
		CourseId string `json:"courseId"`
	}

	var classCourseIds ClassCourseIds

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&classCourseIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	classHandler.log.Println(classCourseIds)

	updatedClass, err := classHandler.classService.AddCourseToClass(token, classCourseIds.ClassId, classCourseIds.CourseId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, updatedClass)
}

func (classHandler RestClassHandlerAdapter) GetRemoveClassCourse(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type ClassCourseIds struct {
		ClassId  string `json:"classId"`
		CourseId string `json:"courseId"`
	}

	var classCourseIds ClassCourseIds

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&classCourseIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	updatedClass, err := classHandler.classService.RemoveCourseFromClass(token, classCourseIds.ClassId, classCourseIds.CourseId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, updatedClass)
}

func (classHandler RestClassHandlerAdapter) GetClassMasterSheet(w http.ResponseWriter, r *http.Request) {
	// token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// var classId string

	// classId = r.URL.Query()["id"][0]

	link, err := classHandler.classService.GetClassMasterSheet("token", "classId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// classHandler.log.Println(link)
	class_utils.SendMarshaledResponse(w, r.Host+link)
	// http.Redirect(w, r, r.Host+link, http.StatusTemporaryRedirect)
}
