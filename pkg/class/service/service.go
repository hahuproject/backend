package class_service

import (
	"log"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/repo"
)

type ClassServicePort interface {

	//Course
	AddCourse(token string, course class_domain.Course) (class_domain.Course, error)
	GetCourses(token string) ([]class_domain.Course, error)
	GetCourse(token, id string) (class_domain.Course, error)
	UpdateCourse(token string, course class_domain.Course) (class_domain.Course, error)
	DeleteCourse(token, id string) error
	AddUserToCourse(token, courseId, userId string) (class_domain.Course, error)
	RemoveUserFromCourse(token, courseId, userId string) (class_domain.Course, error)

	//Department
	AddDepartment(token string, department class_domain.Department) (class_domain.Department, error)
	GetDepartments(token string) ([]class_domain.Department, error)
	GetDepartment(token string, id string) (class_domain.Department, error)
	UpdateDepartment(token string, department class_domain.Department) (class_domain.Department, error)
	DeleteDepartment(token, departmentId string) error

	//Class
	AddClass(token string, class class_domain.Class) (class_domain.Class, error)
	GetClasses(token string) ([]class_domain.Class, error)
	GetClass(token, id string) (class_domain.Class, error)
	UpdateClass(token string, class class_domain.Class) (class_domain.Class, error)
	DeleteClass(token, id string) error
	AddCourseToClass(token, classId, courseId string) (class_domain.Class, error)
	RemoveCourseFromClass(token, classId, courseId string) (class_domain.Class, error)
	GetClassMasterSheet(token, classId string) (string, error)

	//Section
	AddSection(token string, section class_domain.Section) (class_domain.Section, error)
	GetSections(token string) ([]class_domain.Section, error)
	GetSection(token, sectionId string) (class_domain.Section, error)
	AddStudentToSection(token, userId, sectionId string, courses []string) (class_domain.Section, error)
	RemoveStudentFromSection(token, userId, sectionId string, courses []string) (class_domain.Section, error)
	UpdateSection(token string, section class_domain.Section) (class_domain.Section, error)
	DeleteSection(token, id string) error

	//Session
	AddSession(token string, session class_domain.Session) (class_domain.Session, error)
	GetSessions(token string) ([]class_domain.Session, error)
	DeleteSession(token, sessionId string) error

	//Stream
	AddStream(token string, stream class_domain.Stream) (class_domain.Stream, error)
	// GetSessions(token string) ([]class_domain.Session, error)
	UpdateStream(token string, stream class_domain.Stream) (class_domain.Stream, error)
	DeleteStream(token, streamId string) error
}

type ClassServiceAdapter struct {
	repo class_repo.ClassRepoPort
	log  *log.Logger
}

func NewClassServiceAdapter(log *log.Logger, repo class_repo.ClassRepoPort) ClassServicePort {
	return &ClassServiceAdapter{log: log, repo: repo}
}
