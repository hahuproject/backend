package message_repo

import (
	message_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/domain"
)

type MessageRepoPort interface {
	StoreMessage(message message_domain.Message) (message_domain.Message, error)
	FindUserMessages(id string) ([]message_domain.Message, error)
	ReadMessage(messageId string, userId string) (message_domain.Message, error)
}
