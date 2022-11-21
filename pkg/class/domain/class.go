package class_domain

import auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"

type Class struct {
	ID         string                `json:"classId"`
	Name       string                `json:"name"`
	SubName    string                `json:"subName"`
	Students   []auth_domain.Student `json:"students"`
	Department Department            `json:"department"`
	Stream     Stream                `json:"stream"`
	Courses    []Course              `json:"courses"`
	Sections   []Section             `json:"sections"`
}
