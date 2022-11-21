package announcement_rest_handler

import (
	"encoding/json"
	"log"
	"net/http"

	annoucement_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/domain"
	announcement_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/handler"
	annoucement_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/service"
	announcement_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/utils"
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	"golang.org/x/net/websocket"
)

type RestAnnouncementHanlerAdapter struct {
	log     *log.Logger
	service annoucement_service.AnnouncementServicePort
	clients map[string]*websocket.Conn
}

func NewRestAnnouncementHanlerAdapter(log *log.Logger, service annoucement_service.AnnouncementServicePort) announcement_handler.RestAnnouncementHanlerPort {
	return &RestAnnouncementHanlerAdapter{log: log, service: service, clients: make(map[string]*websocket.Conn)}
}

func (handler RestAnnouncementHanlerAdapter) GetConnectToAnnouncementSocket(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(func(ws *websocket.Conn) {

		type AnnouncementWSMessageType struct {
			Event string      `json:"event"`
			Data  interface{} `json:"data"`
		}

		var user auth_domain.User

		var token string

		websocket.JSON.Send(ws, AnnouncementWSMessageType{Event: "STATE", Data: "Connecting to announcement socket server"})
		token = ws.Request().Header.Get("Authorization")

		if token == "" {
			token = r.URL.Query()["token"][0]
		}

		if token == "" {
			log.Println("Message User Empty Token Error")
			ws.Close()
			return
		}

		user, err := announcement_utils.CheckAuth(token)
		if err != nil {
			// log.Println("Message User Auth Req Error", err)
			ws.Close()
			return
		}

		websocket.JSON.Send(ws, AnnouncementWSMessageType{Event: "STATE", Data: "Connected to messaging server"})

		handler.clients[user.ID] = ws

		// var msg annoucement_domain.Announcement
		for {

			var _recievedMsg AnnouncementWSMessageType

			err := websocket.JSON.Receive(ws, &_recievedMsg)
			if err != nil {
				log.Println(err)
				return
			}

			//Broadcast message
			keys := make([]string, 0, len(handler.clients))
			for k := range handler.clients {
				if k != user.ID {
					keys = append(keys, k)
				}
			}

			for i := 0; i < len(keys); i++ {
				err = websocket.JSON.Send(handler.clients[keys[i]], _recievedMsg)
				if err != nil {
					log.Println(err)
					// return
				}
			}

		}
	})}
	s.ServeHTTP(w, r)
}
func (handler RestAnnouncementHanlerAdapter) GetAddAnnouncement(w http.ResponseWriter, r *http.Request) {
	token, err := announcement_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var announcement annoucement_domain.Announcement

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&announcement)
	if err != nil {
		// handler.log.Println(err)
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	addedAnnouncement, err := handler.service.AddAnnouncement(token, announcement)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	announcement_utils.SendMarshaledResponse(w, addedAnnouncement)

}
func (handler RestAnnouncementHanlerAdapter) GetUpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	token, err := announcement_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var announcement annoucement_domain.Announcement

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&announcement)
	if err != nil {
		handler.log.Println(err)
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	addedAnnouncement, err := handler.service.UpdateAnnouncement(token, announcement)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	announcement_utils.SendMarshaledResponse(w, addedAnnouncement)
}

func (handler RestAnnouncementHanlerAdapter) GetDeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	token, err := announcement_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	announcementId := r.URL.Query().Get("id")

	if announcementId == "" {
		announcement_utils.SendResponse(w, http.StatusBadRequest, "Empty announcement id")
		return
	}

	err = handler.service.DeleteAnnouncement(token, announcementId)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	announcement_utils.SendResponse(w, http.StatusOK, "Successfully deleted announcement")
}

func (handler RestAnnouncementHanlerAdapter) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	token, err := announcement_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var announcements []annoucement_domain.Announcement = make([]annoucement_domain.Announcement, 0)

	sectionId := r.URL.Query().Get("section")

	if sectionId != "" {
		announcements, err = handler.service.GetAnnouncementsBySection(token, sectionId)
	} else {
		announcements, err = handler.service.GetAnnouncements(token)
	}

	if err != nil {
		announcement_utils.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	announcement_utils.SendMarshaledResponse(w, announcements)

}
