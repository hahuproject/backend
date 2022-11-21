package message_domain

import (
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

type Message struct {
	ID        string           `json:"messageId"`
	From      auth_domain.User `json:"from"`
	To        auth_domain.User `json:"to"`
	Content   string           `json:"content"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	Read      bool             `json:"read"`
}
