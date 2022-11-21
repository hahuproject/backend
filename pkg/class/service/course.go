package class_service

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (courseService ClassServiceAdapter) AddCourse(token string, course class_domain.Course) (class_domain.Course, error) {
	var addedCourse class_domain.Course
	var user auth_domain.User

	//Check token
	if token == "" {
		return addedCourse, class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedCourse, class_utils.ErrUnauthorized
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" && user.Type != "DEPARTMENT_HEAD" {
		return addedCourse, class_utils.ErrUnauthorized
	}

	//Check Data
	if course.Name == "" {
		return addedCourse, class_utils.ErrInvalidCourseName
	}
	// if course.CreditHr <= 0 {
	// 	return addedCourse, class_utils.ErrInvalidCourseCreditHr
	// }
	if course.Color == "" {
		return addedCourse, class_utils.ErrInvalidCourseColor
	}

	return courseService.repo.StoreCourse(course)
}
func (courseService ClassServiceAdapter) GetCourses(token string) ([]class_domain.Course, error) {
	var user auth_domain.User
	var courses []class_domain.Course

	//Check token
	if token == "" {
		return courses, class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return courses, class_utils.ErrUnauthorized
	}

	if user.ID == "" {
		return courses, class_utils.ErrUnauthorized
	}

	return courseService.repo.FindCourses()
}
func (courseService ClassServiceAdapter) GetCourse(token string, id string) (class_domain.Course, error) {
	// var user auth_domain.User
	var course class_domain.Course

	//Check token
	if token == "" {
		return course, class_utils.ErrUnauthorized
	}

	_, err := class_utils.CheckAuth(token)
	if err != nil {
		return course, class_utils.ErrUnauthorized
	}

	return courseService.repo.FindCourse(id)
}
func (courseService ClassServiceAdapter) UpdateCourse(token string, course class_domain.Course) (class_domain.Course, error) {
	var user auth_domain.User
	var updatedCourse class_domain.Course

	//Check token
	if token == "" {
		return updatedCourse, class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedCourse, class_utils.ErrUnauthorized
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" && user.Type != "DEPARTMENT_HEAD" {
		return updatedCourse, class_utils.ErrUnauthorized
	}
	return courseService.repo.UpdateCourse(course)
}
func (courseService ClassServiceAdapter) DeleteCourse(token string, id string) error {
	var user auth_domain.User

	//Check token
	if token == "" {
		return class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return class_utils.ErrUnauthorized
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" && user.Type != "DEPARTMENT_HEAD" {
		return class_utils.ErrUnauthorized
	}
	return courseService.repo.DeleteCourse(id)
}
func (courseService ClassServiceAdapter) AddUserToCourse(token string, courseId, userId string) (class_domain.Course, error) {
	var user auth_domain.User
	var updatedCourse class_domain.Course

	//Check token
	if token == "" {
		return updatedCourse, class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedCourse, class_utils.ErrUnauthorized
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" && user.Type != "DEPARTMENT_HEAD" {
		return updatedCourse, class_utils.ErrUnauthorized
	}
	return courseService.repo.AddUserToCourse(courseId, userId)
}
func (courseService ClassServiceAdapter) RemoveUserFromCourse(token string, courseId, userId string) (class_domain.Course, error) {
	var user auth_domain.User

	//Check token
	if token == "" {
		return class_domain.Course{}, class_utils.ErrUnauthorized
	}

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return class_domain.Course{}, class_utils.ErrUnauthorized
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" && user.Type != "DEPARTMENT_HEAD" {
		return class_domain.Course{}, class_utils.ErrUnauthorized
	}

	return courseService.repo.RemoveUserFromCourse(courseId, userId)
}
