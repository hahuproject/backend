package message_rest_handler

import (
	"encoding/json"
	"log"
	"net/http"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	message_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/domain"
	message_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/handler"
	message_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/service"
	message_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/utils"
	"golang.org/x/net/websocket"
)

type MessageRestHandlerAdapter struct {
	log            *log.Logger
	messageService message_service.MessageServicePort
	clients        map[string]*websocket.Conn
}

func NewMessageRestHandlerAdapter(log *log.Logger, messageService message_service.MessageServicePort) message_handler.MessageHandlerPort {
	return &MessageRestHandlerAdapter{log: log, messageService: messageService, clients: make(map[string]*websocket.Conn)}
}

func (messageHandler MessageRestHandlerAdapter) GetConnectToMessageService(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(func(ws *websocket.Conn) {

		messageHandler.log.Println("Connect")
		type MessageWSMessageType struct {
			Event string      `json:"event"`
			Data  interface{} `json:"data"`
		}

		var user auth_domain.User

		var token string

		websocket.JSON.Send(ws, MessageWSMessageType{Event: "CONNECTING", Data: "Connecting to messaging socket server"})
		token = ws.Request().Header.Get("Authorization")

		if token == "" {
			token = r.URL.Query()["token"][0]
		}

		if token == "" {
			log.Println("Message User Empty Token Error")
			ws.Close()
			return
		}

		user, err := message_utils.CheckAuth(token)
		if err != nil {
			log.Println("Message User Auth Req Error", err)
			ws.Close()
			return
		}

		websocket.JSON.Send(ws, MessageWSMessageType{Event: "CONNECTED", Data: "Connected to messaging server"})

		messageHandler.clients[user.ID] = ws

		messageHandler.log.Println("connected")
		messageHandler.log.Println(user.ID)
		for {
			var _recievedMsg MessageWSMessageType
			var _msg message_domain.Message

			messageHandler.log.Println("msg to recieve")

			err := websocket.JSON.Receive(ws, &_recievedMsg)
			if err != nil {
				messageHandler.log.Println("msg to recieve 2")
				log.Println(err)
				return
			}

			messageHandler.log.Println("msg to recieve 1")
			err = json.Unmarshal([]byte(_recievedMsg.Data.(string)), &_msg)
			if err != nil {
				messageHandler.log.Println(err)
			}

			_msg.From = user
			messageHandler.log.Println(_msg)
			// _msg = _recievedMsg.Data.(message_domain.Message)
			// messageHandler.log.Println("msg to recieve 1-1")
			// messageHandler.log.Println("msg to recieve 3")

			// // err = json.Unmarshal([]byte(_recievedMsg.Data), &_msg)
			// messageHandler.log.Println("msg to recieve 4")
			// if err != nil {
			// 	messageHandler.log.Println("msg to recieve 5")
			// 	messageHandler.log.Println(err)
			// }
			// messageHandler.log.Println("msg to recieve 5")
			if messageHandler.clients[_msg.To.ID] != nil {
				err = websocket.JSON.Send(messageHandler.clients[_msg.To.ID], MessageWSMessageType{Event: "RECIEVE", Data: _msg})
				if err != nil {
					log.Println(err)
					return
				}
			}

		}
	})}
	s.ServeHTTP(w, r)
}
func (messageHandler MessageRestHandlerAdapter) GetUserMessages(w http.ResponseWriter, r *http.Request) {

	// messageHandler.log.Println("get user messages 0")
	token, err := message_utils.CheckBearerTokenFromHTTPRequest(r)
	// messageHandler.log.Println("get user messages 1")

	if err != nil {
		// messageHandler.log.Println("get user messages 2")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// messageHandler.log.Println("get user messages 3")

	userMessages, err := messageHandler.messageService.GetUserMessages(token)
	// messageHandler.log.Println("get user messages 4")
	if err != nil {
		// messageHandler.log.Println("get user messages 5")
		// messageHandler.log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// messageHandler.log.Println("get user messages 6")
	message_utils.SendMarshaledResponse(w, userMessages)
	// messageHandler.log.Println("get user messages 7")

}

func (messageHandler MessageRestHandlerAdapter) GetAddMessage(w http.ResponseWriter, r *http.Request) {

	token, err := message_utils.CheckBearerTokenFromHTTPRequest(r)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var message message_domain.Message

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&message)
	if err != nil {
		log.Println("Message User Empty Token Error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode body"))
		return
	}

	defer r.Body.Close()

	userMessages, err := messageHandler.messageService.AddMessage(token, message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	message_utils.SendMarshaledResponse(w, userMessages)

}

func (messageHandler MessageRestHandlerAdapter) GetReadMessage(w http.ResponseWriter, r *http.Request) {

	token, err := message_utils.CheckBearerTokenFromHTTPRequest(r)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var messageId = r.URL.Query().Get("id")

	defer r.Body.Close()

	userMessages, err := messageHandler.messageService.ReadMessage(token, messageId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	message_utils.SendMarshaledResponse(w, userMessages)

}
