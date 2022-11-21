package attendance_rest_handler

import (
	"encoding/json"
	"log"
	"net/http"

	attendance_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/domain"
	attendance_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/handler"
	attendance_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/service"
	attendance_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/utils"
)

type AttendanceRestHandlerAdapter struct {
	log     *log.Logger
	service attendance_service.AttendanceServicePort
}

func NewAttendanceRestHandlerAdapter(log *log.Logger, service attendance_service.AttendanceServicePort) attendance_handler.AttendanceHandlerPort {
	return AttendanceRestHandlerAdapter{log: log, service: service}
}

func (handler AttendanceRestHandlerAdapter) GetAddAttendance(w http.ResponseWriter, r *http.Request) {
	token, err := attendance_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var attendances []attendance_domain.Attendance

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&attendances)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedAttendance, err := handler.service.AddAttendance(token, attendances)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	attendance_utils.SendMarshaledResponse(w, addedAttendance)

}
func (handler AttendanceRestHandlerAdapter) GetAttendances(w http.ResponseWriter, r *http.Request) {
	var attendances []attendance_domain.Attendance
	var err error

	token, err := attendance_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	userId := r.URL.Query().Get("user-id")

	if userId != "" {
		attendances, err = handler.service.GetUserAttendances(token, userId)
	} else {
		attendances, err = handler.service.GetAttendances(token)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	attendance_utils.SendMarshaledResponse(w, attendances)

}
