package class_rest_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (handler RestClassHandlerAdapter) GetAddStream(w http.ResponseWriter, r *http.Request) {
	token, err := class_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var stream class_domain.Stream

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&stream)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	addedStream, err := handler.classService.AddStream(token, stream)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	class_utils.SendMarshaledResponse(w, addedStream)
}

func (handler RestClassHandlerAdapter) GetUpdateStream(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	var stream class_domain.Stream

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&stream)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error occured while decoding body"))
		return
	}
	defer r.Body.Close()

	updatedStream, err := handler.classService.UpdateStream(token, stream)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUpdatedCourse, _ := json.Marshal(updatedStream)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUpdatedCourse)
}

func (handler RestClassHandlerAdapter) GetDeleteStream(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || len(strings.Split(token, " ")) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized request"))
		return
	}

	token = strings.Split(token, " ")[1]

	err := handler.classService.DeleteStream(token, r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted stream"))
}

// func (sessionHandler RestClassHandlerAdapter) GetStreams(w http.ResponseWriter, r *http.Request)   {}
// func (sessionHandler RestClassHandlerAdapter) GetStream(w http.ResponseWriter, r *http.Request)    {}
// func (sessionHandler RestClassHandlerAdapter) GetUpdateStream(w http.ResponseWriter, r *http.Request) {
// }
// func (sessionHandler RestClassHandlerAdapter) GetDeleteStream(w http.ResponseWriter, r *http.Request) {
// }
// func (sessionHandler RestClassHandlerAdapter) GetStreamsByDepartment(w http.ResponseWriter, r *http.Request) {
// }
