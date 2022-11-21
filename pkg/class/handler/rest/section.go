package class_rest_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (sectionHandler RestClassHandlerAdapter) GetAddSection(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var section class_domain.Section

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&section)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedSection, err := sectionHandler.classService.AddSection(token, section)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledAddedSection, err := json.Marshal(addedSection)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedSection)
}
func (sectionHandler RestClassHandlerAdapter) GetSections(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sectionHandler.log.Println("get sections 0")
	sections, err := sectionHandler.classService.GetSections(token)
	if err != nil {
		sectionHandler.log.Println("get sections 01")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sectionHandler.log.Println("get sections 1")

	marshaledSections, err := json.Marshal(sections)
	if err != nil {
		sectionHandler.log.Println("get sections 2")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sectionHandler.log.Println("get sections 3")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledSections)
}
func (sectionHandler RestClassHandlerAdapter) GetSection(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(r.URL.Query()["id"]) < 1 || r.URL.Query()["id"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid section id"))
		return
	}

	section, err := sectionHandler.classService.GetSection(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledSection, err := json.Marshal(section)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledSection)
}

func (sectionHandler RestClassHandlerAdapter) GetAddStudentToSection(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type userSectionIDs struct {
		UserId    string   `json:"userId"`
		SectionId string   `json:"sectionId"`
		Courses   []string `json:"courses"`
	}

	var userSectionId userSectionIDs

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&userSectionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	updatedSection, err := sectionHandler.classService.AddStudentToSection(token, userSectionId.UserId, userSectionId.SectionId, userSectionId.Courses)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	class_utils.SendMarshaledResponse(w, updatedSection)
}
func (sectionHandler RestClassHandlerAdapter) GetRemoveStudentFromSection(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	type userSectionIDs struct {
		UserId    string   `json:"userId"`
		SectionId string   `json:"sectionId"`
		Courses   []string `json:"courses"`
	}

	var userSectionId userSectionIDs

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&userSectionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	updatedSection, err := sectionHandler.classService.RemoveStudentFromSection(token, userSectionId.UserId, userSectionId.SectionId, userSectionId.Courses)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, updatedSection)
}

func (classHandler RestClassHandlerAdapter) GetUpdateSection(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	var section class_domain.Section

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&section)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured while decoding body"))
		return
	}
	defer r.Body.Close()

	updatedSection, err := classHandler.classService.UpdateSection(token, section)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUpdatedCourse, _ := json.Marshal(updatedSection)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUpdatedCourse)
}
func (classHandler RestClassHandlerAdapter) GetDeleteSection(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	err := classHandler.classService.DeleteSection(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted section"))
}
