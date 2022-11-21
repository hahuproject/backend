package attendance_repo

import (
	attendance_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/domain"
)

type AttendanceRepoPort interface {
	StoreAttendance(attendance []attendance_domain.Attendance, userId string) ([]attendance_domain.Attendance, error)
	// FindAttendances() ([]attendance_domain.Attendance, error)
	FindAttendancesByUser(userId string) ([]attendance_domain.Attendance, error)
	FindAttendancesByInstructor(instructorId string) ([]attendance_domain.Attendance, error)
	FindAttendancesBySession(sessionId string) ([]attendance_domain.Attendance, error)
}
