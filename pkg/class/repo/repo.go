package class_repo

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type ClassRepoPort interface {

	//Course

	StoreCourse(course class_domain.Course) (class_domain.Course, error)
	FindCourses() ([]class_domain.Course, error)
	FindCourse(id string) (class_domain.Course, error)
	UpdateCourse(course class_domain.Course) (class_domain.Course, error)
	DeleteCourse(id string) error
	AddUserToCourse(courseId, userId string) (class_domain.Course, error)
	RemoveUserFromCourse(courseId, userId string) (class_domain.Course, error)

	//Department

	StoreDepartment(department class_domain.Department) (class_domain.Department, error)
	FindDepartments() ([]class_domain.Department, error)
	FindDepartmentsByHead(head auth_domain.User) ([]class_domain.Department, error)
	FindDepartment(id string) (class_domain.Department, error)
	UpdateDepartment(department class_domain.Department) (class_domain.Department, error)
	DeleteDepartment(departmentId string) error

	//Class
	StoreClass(class class_domain.Class) (class_domain.Class, error)
	FindClasses() ([]class_domain.Class, error)
	FindClass(id string) (class_domain.Class, error)
	UpdateClass(class class_domain.Class) (class_domain.Class, error)
	DeleteClass(id string) error
	FindClassesByDepartmentHead(id string) ([]class_domain.Class, error)
	StoreClassCourse(classId, courseId string) (class_domain.Class, error)
	RemoveClassCourse(classId, courseId string) (class_domain.Class, error)

	//Section
	StoreSection(section class_domain.Section) (class_domain.Section, error)
	FindSections() ([]class_domain.Section, error)
	FindSection(sectionId string) (class_domain.Section, error)
	FindSectionsByDepartmentHead(departmentHeadId string) ([]class_domain.Section, error)
	FindSectionsByInstructor(instructorId string) ([]class_domain.Section, error)
	UpdateSection(section class_domain.Section) (class_domain.Section, error)
	DeleteSection(id string) error
	AddStudentToSection(userId, sectionId string, courses []string) (class_domain.Section, error)
	RemoveStudentFromSection(userId, sectionId string, courses []string) (class_domain.Section, error)

	//Session
	StoreSession(session class_domain.Session) (class_domain.Session, error)
	FindSessions() ([]class_domain.Session, error)
	FindSessionBySection(sectionId string) ([]class_domain.Session, error)
	FindSessionsByStudent(sectionId string) ([]class_domain.Session, error)
	FindSessionsByInstructor(instructorId string) ([]class_domain.Session, error)
	DeleteSession(sessionId string) error

	//Stream
	StoreStream(stream class_domain.Stream) (class_domain.Stream, error)
	UpdateStream(stream class_domain.Stream) (class_domain.Stream, error)
	DeleteStream(streamId string) error
}
