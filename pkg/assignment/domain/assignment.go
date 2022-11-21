package assignment_domain

import (
	"time"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Assignment struct {
	ID         string               `json:"assignmentId"`
	Title      string               `json:"title"`
	Remark     string               `json:"remark"`
	Attachment string               `json:"attachment"`
	CreatedAt  time.Time            `json:"createdAt"`
	Section    class_domain.Section `json:"section"`
	Course     class_domain.Course  `json:"course"`
}
