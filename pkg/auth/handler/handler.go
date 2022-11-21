package auth_handler

import "net/http"

type RestAuthHandlerPort interface {
	GetRegisterAdmin(w http.ResponseWriter, r *http.Request)
	GetRegisterRegistryOfficer(w http.ResponseWriter, r *http.Request)
	GetRegisterSubRegistryOfficer(w http.ResponseWriter, r *http.Request)
	GetRegisterDepartmentHead(w http.ResponseWriter, r *http.Request)
	GetRegisterInstructor(w http.ResponseWriter, r *http.Request)
	GetRegisterSupervisor(w http.ResponseWriter, r *http.Request)
	GetRegisterStudent(w http.ResponseWriter, r *http.Request)
	GetLogin(w http.ResponseWriter, r *http.Request)
	GetApproveUser(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetChangePassword(w http.ResponseWriter, r *http.Request)
	GetUpdateProfile(w http.ResponseWriter, r *http.Request)
	GetBanUser(w http.ResponseWriter, r *http.Request)
	GetUnBanUser(w http.ResponseWriter, r *http.Request)

	//GetStudents
	GetStudents(w http.ResponseWriter, r *http.Request)
	GetUpdateStudentPayment(w http.ResponseWriter, r *http.Request)
}
