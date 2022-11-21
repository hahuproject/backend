package grade_domain

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Grade struct {
	ID              string               `json:"gradeId"`
	Assessment      float64              `json:"assessment"`
	Mid             float64              `json:"mid"`
	Final           float64              `json:"final"`
	Lab             float64              `json:"lab"`
	User            auth_domain.User     `json:"user"`
	Course          class_domain.Course  `json:"course"`
	CanReview       bool                 `json:"canReview"`
	ReviewRequested bool                 `json:"reviewRequested"`
	Sunmitted       bool                 `json:"submitted"`
	Section         class_domain.Section `json:"section"`
}
