package assignment_rest_handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	assignment_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/domain"
	assignment_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/handler"
	assignment_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/service"
	assignment_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/utils"
)

type AssignmentRestHandlerAdapter struct {
	log     *log.Logger
	service assignment_service.AssignmentServiceRepo
}

func NewAssignmentRestHandlerAdapter(log *log.Logger, service assignment_service.AssignmentServiceRepo) assignment_handler.AssignmentHandlerPort {
	return AssignmentRestHandlerAdapter{log: log, service: service}
}

func (handler AssignmentRestHandlerAdapter) GetAssignments(w http.ResponseWriter, r *http.Request) {
	token, err := assignment_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	assignments, err := handler.service.GetAssignments(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	assignment_utils.SendMarshaledResponse(w, assignments)
}
func (handler AssignmentRestHandlerAdapter) GetAddAssignment(w http.ResponseWriter, r *http.Request) {
	token, err := assignment_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var assignment assignment_domain.Assignment

	//Get Assignment File
	assignmentFile, assignmentFileHandler, err := r.FormFile("attachment")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a valid attachment"))
		return
	}

	defer assignmentFile.Close()

	_, err = os.Open("/uploads/files/assignment-attachments")
	if err != nil {
		err = os.MkdirAll("uploads/files/assignment-attachments", 0755)
		if err != nil {
			handler.log.Println(err)
		}
	}

	// _, err = os.Open("/uploads/images/book-cover")
	// if err != nil {
	// 	err = os.MkdirAll("uploads/images/book-cover", 0755)
	// 	if err != nil {
	// 		libraryHandler.log.Println(err)
	// 	}
	// }

	storedAssignmentFile, err := os.OpenFile("./uploads/files/assignment-attachments/"+r.FormValue("title")+strings.Split(assignmentFileHandler.Filename, ".")[len(strings.Split(assignmentFileHandler.Filename, "."))-1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		handler.log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to save attachment"))
		return
	}

	io.Copy(storedAssignmentFile, assignmentFile)

	assignment.Attachment = r.Host + strings.Replace(storedAssignmentFile.Name(), ".", "", 1)

	assignment.Title = r.FormValue("title")
	assignment.Remark = r.FormValue("remark")
	assignment.Section.ID = r.FormValue("section")
	assignment.Course.ID = r.FormValue("course")

	addedAssignment, err := handler.service.AddAssignment(token, assignment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	assignment_utils.SendMarshaledResponse(w, addedAssignment)

}
func (handler AssignmentRestHandlerAdapter) GetUpdateAssignment(w http.ResponseWriter, r *http.Request) {
	token, err := assignment_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var assignment assignment_domain.Assignment

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&assignment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	updatedAssignment, err := handler.service.UpdateAssignment(token, assignment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	assignment_utils.SendMarshaledResponse(w, updatedAssignment)

}
func (handler AssignmentRestHandlerAdapter) GetDeleteAssignment(w http.ResponseWriter, r *http.Request) {
	token, err := assignment_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	assignmentId := r.URL.Query().Get("id")

	if assignmentId == "" {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No assignment id specified"))
			return
		}
	}

	err = handler.service.DeleteAssignment(token, assignmentId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted assignment"))

}
