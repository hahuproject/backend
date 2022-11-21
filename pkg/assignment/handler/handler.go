package assignment_handler

import "net/http"

type AssignmentHandlerPort interface {
	GetAssignments(w http.ResponseWriter, r *http.Request)
	GetAddAssignment(w http.ResponseWriter, r *http.Request)
	GetUpdateAssignment(w http.ResponseWriter, r *http.Request)
	GetDeleteAssignment(w http.ResponseWriter, r *http.Request)
}
