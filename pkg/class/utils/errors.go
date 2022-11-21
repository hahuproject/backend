package class_utils

import "errors"

var (
	//AUTH
	ErrUnauthorized = errors.New("unautorized request")

	//REPO - COURSE
	ErrCourseNameExists = errors.New("course name is taken")

	//SERVICE - COURSE
	ErrInvalidCourseName     = errors.New("invalid course name")
	ErrInvalidCourseCreditHr = errors.New("invalid course credit hour")
	ErrInvalidCourseColor    = errors.New("invalid course color")

	//SERVICE - DEPARTMENT
	ErrInvalidDepartmentName = errors.New("invalid department name")

	ErrDbPingFailed = errors.New("failed to ping db")
)
