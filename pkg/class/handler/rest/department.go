package class_rest_handler

import (
	"encoding/json"
	"net/http"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (departmentHandler RestClassHandlerAdapter) GetAddDepartment(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var department class_domain.Department

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&department)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedDepartment, err := departmentHandler.classService.AddDepartment(token, department)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledAddedDepartment, err := json.Marshal(addedDepartment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedDepartment)
}
func (departmentHandler RestClassHandlerAdapter) GetDepartments(w http.ResponseWriter, r *http.Request) {
	var token string
	token, _ = class_utils.CheckBearerTokenFromHTTPRequest(r)

	departments, err := departmentHandler.classService.GetDepartments(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, departments)
}
func (departmentHandler RestClassHandlerAdapter) GetDepartment(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(r.URL.Query()["id"]) < 1 || r.URL.Query()["id"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid department id"))
		return
	}

	department, err := departmentHandler.classService.GetDepartment(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledDepartment, err := json.Marshal(department)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledDepartment)
}

func (departmentHandler RestClassHandlerAdapter) GetUpdateDepartment(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var department class_domain.Department

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&department)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	departmentHandler.log.Println(department.Head.ID)

	updatedDepartment, err := departmentHandler.classService.UpdateDepartment(token, department)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUpdatedDepartment, _ := json.Marshal(updatedDepartment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUpdatedDepartment)
}

func (departmentHandler RestClassHandlerAdapter) GetDeleteDepartment(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	departmentId := r.URL.Query().Get("id")
	if departmentId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Specify department id"))
		return
	}

	err = departmentHandler.classService.DeleteDepartment(token, departmentId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted department"))
}
