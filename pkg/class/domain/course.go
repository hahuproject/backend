package class_domain

import auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"

type Course struct {
	ID            string             `json:"courseId"`
	Name          string             `json:"name"`
	CreditHr      int                `json:"creditHr"`
	Color         string             `json:"color"`
	Users         []auth_domain.User `json:"users"`
	Prerequisites []Course           `json:"prerequisites"`
}
