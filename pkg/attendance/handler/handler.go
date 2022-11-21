package attendance_handler

import "net/http"

type AttendanceHandlerPort interface {
	GetAddAttendance(w http.ResponseWriter, r *http.Request)
	GetAttendances(w http.ResponseWriter, r *http.Request)
}
