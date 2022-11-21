package auth_repo

import (
	auth "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

type AuthRepositoryPort interface {
	StoreAdmin(user auth.Admin) (auth.Admin, error)
	StoreRegistryOfficer(user auth.RegistryOfficer) (auth.RegistryOfficer, error)
	StoreSubRegistryOfficer(user auth.SubRegistryOfficer) (auth.SubRegistryOfficer, error)
	StoreDepartmentHead(user auth.DepartmentHead) (auth.DepartmentHead, error)
	StoreInstructor(user auth.Instructor) (auth.Instructor, error)
	StoreSupervisor(user auth.Supervisor) (auth.Supervisor, error)
	StoreStudent(verified bool, user auth.Student) (auth.Student, error)
	FindUserByUsername(username string) (auth.User, error)
	FindUserById(id string) (auth.User, error)
	GetUserRoleByID(id string) (string, error)
	ApproveUser(id string) error
	BanUser(id string) error
	UnBanUser(id string) error
	FindAllUsers() ([]auth.User, error)
	FindStudent(id string) (auth.Student, error)
	FindInstructor(id string) (auth.Instructor, error)
	FindSupervisor(id string) (auth.Supervisor, error)
	ChangePassword(userId, newPassword string) error
	UpdateProfile(user auth.User) error

	//Student
	FindStudents() ([]auth.Student, error)
	UpdateStudentPaymentSatus(userId string, staus bool) (auth.Student, error)
}
