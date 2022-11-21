package announcement_handler

import "net/http"

type RestAnnouncementHanlerPort interface {
	GetConnectToAnnouncementSocket(w http.ResponseWriter, r *http.Request)
	GetAnnouncements(w http.ResponseWriter, r *http.Request)
	GetAddAnnouncement(w http.ResponseWriter, r *http.Request)
	GetUpdateAnnouncement(w http.ResponseWriter, r *http.Request)
	GetDeleteAnnouncement(w http.ResponseWriter, r *http.Request)
}
