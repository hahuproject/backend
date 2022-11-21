package class_domain

import auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"

type Section struct {
	ID       string                `json:"sectionId"`
	Name     string                `json:"name"`
	Year     int64                 `json:"year"`
	Class    Class                 `json:"class"`
	Students []auth_domain.Student `json:"students"`
	Sessions []Session             `json:"sessions"`
}
