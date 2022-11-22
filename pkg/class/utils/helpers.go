package class_utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

var (
	AUTH_URL_DEV       = "http://localhost:5000/auth/me"
	AUTH_URL_PROD      = "https://hahu-sms.herokuapp.com/auth/me"
	AUTH_URL_LOCAL     = "http://18.134.74.142:5003/auth/me"
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

	log.Println("CHECKING AUTH 0")
	if token == "" {
		log.Println("CHECKING AUTH 1")
		return user, ErrUnauthorized
	}

	log.Println("CHECKING AUTH 2")
	req, err := http.NewRequest(http.MethodGet, AUTH_URL, nil)
	if err != nil {
		log.Println("CHECKING AUTH 3")
		return user, err
	}

	log.Println("CHECKING AUTH 4")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	log.Println("CHECKING AUTH 5")
	res, err := client.Do(req)
	if err != nil {
		log.Println("CHECKING AUTH 6")
		return user, err
	}

	log.Println("CHECKING AUTH 7")
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&user)
	log.Println("CHECKING AUTH 8")
	if err != nil {
		log.Println("CHECKING AUTH 9")
		return user, err
	}

	log.Println("CHECKING AUTH 10")
	return user, nil
}

func SendMarshaledResponse(w http.ResponseWriter, data interface{}) {
	marshaledAddedClass, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedClass)
}
