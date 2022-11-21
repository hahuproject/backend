package class_handler

import "net/http"

type ClassHandlerPort interface {

	//Course Handler
	GetAddCourse(w http.ResponseWriter, r *http.Request)
	GetCourses(w http.ResponseWriter, r *http.Request)
	GetCourse(w http.ResponseWriter, r *http.Request)
	GetUpdateCourse(w http.ResponseWriter, r *http.Request)
	GetDeleteCourse(w http.ResponseWriter, r *http.Request)
	GetAddUserToCourse(w http.ResponseWriter, r *http.Request)
	GetRemoveUserFromCourse(w http.ResponseWriter, r *http.Request)

	//Department
	GetAddDepartment(w http.ResponseWriter, r *http.Request)
	GetDepartments(w http.ResponseWriter, r *http.Request)
	GetDepartment(w http.ResponseWriter, r *http.Request)
	GetUpdateDepartment(w http.ResponseWriter, r *http.Request)
	GetDeleteDepartment(w http.ResponseWriter, r *http.Request)

	//Class
	GetAddClass(w http.ResponseWriter, r *http.Request)
	GetClasses(w http.ResponseWriter, r *http.Request)
	GetClass(w http.ResponseWriter, r *http.Request)
	GetUpdateClass(w http.ResponseWriter, r *http.Request)
	GetDeleteClass(w http.ResponseWriter, r *http.Request)
	GetAddClassCourse(w http.ResponseWriter, r *http.Request)
	GetRemoveClassCourse(w http.ResponseWriter, r *http.Request)
	GetClassMasterSheet(w http.ResponseWriter, r *http.Request)

	//Section
	GetAddSection(w http.ResponseWriter, r *http.Request)
	GetSections(w http.ResponseWriter, r *http.Request)
	GetSection(w http.ResponseWriter, r *http.Request)
	GetAddStudentToSection(w http.ResponseWriter, r *http.Request)
	GetRemoveStudentFromSection(w http.ResponseWriter, r *http.Request)
	GetUpdateSection(w http.ResponseWriter, r *http.Request)
	GetDeleteSection(w http.ResponseWriter, r *http.Request)

	//Session
	GetAddSession(w http.ResponseWriter, r *http.Request)
	GetSessions(w http.ResponseWriter, r *http.Request)
	GetDeleteSession(w http.ResponseWriter, r *http.Request)

	//Stream
	GetAddStream(w http.ResponseWriter, r *http.Request)
	GetUpdateStream(w http.ResponseWriter, r *http.Request)
	GetDeleteStream(w http.ResponseWriter, r *http.Request)
	// GetStreams(w http.ResponseWriter, r *http.Request)
	// GetStream(w http.ResponseWriter, r *http.Request)
	// GetUpdateStream(w http.ResponseWriter, r *http.Request)
	// GetDeleteStream(w http.ResponseWriter, r *http.Request)
	// GetStreamsByDepartment(w http.ResponseWriter, r *http.Request)
}
