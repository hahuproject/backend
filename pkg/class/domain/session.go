package class_domain

import (
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

type Session struct {
	ID         string           `json:"sessionId"`
	Duration   int              `json:"duration"`
	StartDate  time.Time        `json:"startDate"`
	Section    Section          `json:"section"`
	Course     Course           `json:"course"`
	Instructor auth_domain.User `json:"instructor"`
}
