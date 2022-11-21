package class_domain

import auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"

type Department struct {
	ID      string                     `json:"departmentId"`
	Name    string                     `json:"name"`
	Head    auth_domain.DepartmentHead `json:"head"`
	Classes []Class                    `json:"classes"`
	Streams []Stream                   `json:"streams"`
}
