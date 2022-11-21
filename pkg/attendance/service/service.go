package attendance_service

import (
	"log"

	attendance_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/domain"
	attendance_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/repo"
	attendance_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/utils"
)

type AttendanceServicePort interface {
	AddAttendance(token string, attendances []attendance_domain.Attendance) ([]attendance_domain.Attendance, error)
	GetAttendances(token string) ([]attendance_domain.Attendance, error)
	GetUserAttendances(token, userId string) ([]attendance_domain.Attendance, error)
}

type AttendanceServiceAdapter struct {
	log  *log.Logger
	repo attendance_repo.AttendanceRepoPort
}

func NewAttendanceServiceAdapter(log *log.Logger, repo attendance_repo.AttendanceRepoPort) AttendanceServicePort {
	return AttendanceServiceAdapter{log: log, repo: repo}
}

func (service AttendanceServiceAdapter) AddAttendance(token string, attendances []attendance_domain.Attendance) ([]attendance_domain.Attendance, error) {

	var addedAttendances []attendance_domain.Attendance

	user, err := attendance_utils.CheckAuth(token)
	if err != nil {
		return addedAttendances, err
	}

	var _attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)
	for i := 0; i < len(attendances); i++ {
		if attendances[i].User.Type == "STUDENT" {
			_attendances = append(_attendances, attendances[i])
		}
	}

	if user.Type != "INSTRUCTOR" && user.Type != "SUPERVISOR" {
		return addedAttendances, attendance_utils.ErrNotAuthorized
	}

	return service.repo.StoreAttendance(_attendances, user.ID)
}

func (service AttendanceServiceAdapter) GetAttendances(token string) ([]attendance_domain.Attendance, error) {

	var attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)

	user, err := attendance_utils.CheckAuth(token)
	if err != nil {
		return attendances, err
	}

	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" {
		return service.repo.FindAttendancesByInstructor(user.ID)
	}

	if user.Type == "STUDENT" {
		return service.repo.FindAttendancesByUser(user.ID)
	}

	if user.Type == "DEPARTMENT_HEAD" {
		return service.repo.FindAttendancesByUser(user.ID)
	}

	return attendances, attendance_utils.ErrNotAuthorized
}

func (service AttendanceServiceAdapter) GetUserAttendances(token, userId string) ([]attendance_domain.Attendance, error) {

	var attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)

	user, err := attendance_utils.CheckAuth(token)
	if err != nil {
		return attendances, err
	}

	if user.Type == "DEPARTMENT_HEAD" {
		return service.repo.FindAttendancesByUser(userId)
	}

	return attendances, attendance_utils.ErrNotAuthorized
}
