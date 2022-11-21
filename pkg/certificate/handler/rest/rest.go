package certificate_rest_handler

import (
	"log"
	"net/http"

	certificate_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/handler"
	certificate_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/service"
	certificate_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/utils"
)

type RestCertificateHandlerAdapter struct {
	log     *log.Logger
	service certificate_service.CertificateServicePort
}

func NewRestCertificateHandlerAdapter(log *log.Logger, service certificate_service.CertificateServicePort) certificate_handler.CertificateHandlerPort {
	return &RestCertificateHandlerAdapter{log, service}
}

func (handler RestCertificateHandlerAdapter) GetCertificate(w http.ResponseWriter, r *http.Request) {

	token, err := certificate_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unauthorized Request"))
		return
	}

	userId := r.URL.Query().Get("user-id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provide a valid user id"))
		return

	}

	generatedCertificate, err := handler.service.GenerateCertificate(token, userId)
	if err != nil {
		handler.log.Println("generatedCertificate err")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	handler.log.Println(generatedCertificate)
	handler.log.Println("generatedCertificate")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("http://" + r.Host + generatedCertificate.File))
}
