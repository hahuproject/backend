package attendance_domain

import (
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Attendance struct {
	ID        string               `json:"attendanceId"`
	User      auth_domain.User     `json:"user"`
	Status    string               `json:"status"`
	Session   class_domain.Session `json:"session"`
	CreatedAt time.Time            `json:"createdAt"`
}
