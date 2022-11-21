package message_handler

import (
	"net/http"
)

type MessageHandlerPort interface {
	GetConnectToMessageService(w http.ResponseWriter, r *http.Request)
	GetUserMessages(w http.ResponseWriter, r *http.Request)
	GetAddMessage(w http.ResponseWriter, r *http.Request)
	GetReadMessage(w http.ResponseWriter, r *http.Request)
}
