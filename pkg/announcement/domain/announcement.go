package annoucement_domain

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Announcement struct {
	ID        string               `json:"announcementId"`
	Title     string               `json:"title"`
	Message   string               `json:"message"`
	CreatedAt string               `json:"createdAt"`
	Section   class_domain.Section `json:"section"`
	PostedBy  auth_domain.User     `json:"postedBy"`
}
