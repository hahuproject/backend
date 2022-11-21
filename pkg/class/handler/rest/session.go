package class_rest_handler

import (
	"encoding/json"
	"net/http"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (sessionHandler RestClassHandlerAdapter) GetAddSession(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var session class_domain.Session

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&session)
	if err != nil {
		sessionHandler.log.Println("error in decode" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionHandler.log.Println(session.Duration)
	sessionHandler.log.Println(session.StartDate)
	sessionHandler.log.Println(session.Section.ID)
	sessionHandler.log.Println(session.Course.ID)
	sessionHandler.log.Println(session.Instructor.ID)
	defer r.Body.Close()

	addedSession, err := sessionHandler.classService.AddSession(token, session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledAddedSession, err := json.Marshal(addedSession)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedSession)
}

func (sessionHandler RestClassHandlerAdapter) GetSessions(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	studentSessions, err := sessionHandler.classService.GetSessions(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, studentSessions)
}

func (handler RestClassHandlerAdapter) GetDeleteSession(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionId := r.URL.Query().Get("id")
	if sessionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not valid session id"))
		return
	}

	err = handler.classService.DeleteSession(token, sessionId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted session"))
}
