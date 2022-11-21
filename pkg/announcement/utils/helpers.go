package annoucement_utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

var (
	AUTH_URL_DEV       = "http://localhost:5000/auth/me"
	AUTH_URL_PROD      = "https://hahu-sms.herokuapp.com/auth/me"
	AUTH_URL_LOCAL     = "http://192.168.2.222:5002/auth/me"
	AUTH_URL_TEGBAREED = "http://10.3.120.6:5000/auth/me"
	AUTH_URL_VPS       = "http://18.134.74.142:5000/auth/me"
)

var AUTH_URL = AUTH_URL_LOCAL

func CheckBearerTokenFromHTTPRequest(r *http.Request) (string, error) {
	var token string = ""

	if r.Header.Get("Authorization") == "" {
		return token, ErrUnauthorized
	}

	if len(strings.Split(r.Header.Get("Authorization"), " ")) < 2 {
		return token, ErrUnauthorized
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	if token == "" {
		return token, ErrUnauthorized
	}
	return token, nil
}

func CheckAuth(token string) (auth_domain.User, error) {
	var user auth_domain.User

	if token == "" {
		return user, ErrUnauthorized
	}

	req, err := http.NewRequest(http.MethodGet, AUTH_URL, nil)
	if err != nil {
		return user, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return user, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func SendMarshaledResponse(w http.ResponseWriter, data interface{}) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshaledData)
	if err != nil {
		log.Println(err)
	}
}

func SendResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func JsonDecode(from io.ReadCloser, to interface{}) error {
	decoder := json.NewDecoder(from)

	return decoder.Decode(&to)
}
