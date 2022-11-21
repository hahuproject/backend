package certificate_utils

import (
	"encoding/json"
	"errors"
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
	var token string

	if r.Header.Get("Authorization") == "" {
		return token, errors.New("not authorized for the operation")
	}

	if len(strings.Split(r.Header.Get("Authorization"), " ")) < 2 {
		return token, errors.New("not authorized for the operation")
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]

	return token, nil
}

func CheckAuth(token string) (auth_domain.User, error) {
	var user auth_domain.User

	if token == "" {
		return user, errors.New("not authorized for the operation")
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
