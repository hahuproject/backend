package certificate_domain

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Certificate struct {
	User    auth_domain.User     `json:"user"`
	Section class_domain.Section `json:"section"`
	File    string               `json:"file"`
}
