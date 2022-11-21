package message_service

import (
	"log"

	message_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/domain"
	message_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/repo"
	message_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/utils"
)

type MessageServicePort interface {
	AddMessage(token string, message message_domain.Message) (message_domain.Message, error)
	GetUserMessages(token string) ([]message_domain.Message, error)
	ReadMessage(token, messageId string) (message_domain.Message, error)
}

type MessageServiceAdapter struct {
	log  *log.Logger
	repo message_repo.MessageRepoPort
}

func NewMessageServiceAdapter(log *log.Logger, repo message_repo.MessageRepoPort) MessageServicePort {
	return &MessageServiceAdapter{log: log, repo: repo}
}

func (messageService MessageServiceAdapter) AddMessage(token string, message message_domain.Message) (message_domain.Message, error) {

	var userMessages message_domain.Message

	user, err := message_utils.CheckAuth(token)
	if err != nil {
		return userMessages, err
	}

	message.From = user

	return messageService.repo.StoreMessage(message)
}

func (messageService MessageServiceAdapter) GetUserMessages(token string) ([]message_domain.Message, error) {

	var userMessages []message_domain.Message = make([]message_domain.Message, 0)

	// messageService.log.Println("message service 0")
	messageService.log.Println(token)

	if token == "" {
		return userMessages, message_utils.ErrUnauthorized
	}

	// messageService.log.Println("message service 1")

	user, err := message_utils.CheckAuth(token)

	// messageService.log.Println("message service 2")
	if err != nil {
		// messageService.log.Println("message service 3")
		messageService.log.Println(err)
		return userMessages, err
	}
	// messageService.log.Println("message service 4")

	return messageService.repo.FindUserMessages(user.ID)
}

func (messageService MessageServiceAdapter) ReadMessage(token, messageId string) (message_domain.Message, error) {

	var updatedMessage message_domain.Message

	// messageService.log.Println("message service 0")
	user, err := message_utils.CheckAuth(token)

	if err != nil {
		return updatedMessage, err
	}
	messageService.log.Println("message service 4")
	messageService.log.Println(user.ID)
	messageService.log.Println(messageId)

	return messageService.repo.ReadMessage(messageId, user.ID)
}
