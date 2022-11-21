package certificate_handler

import "net/http"

type CertificateHandlerPort interface {
	GetCertificate(w http.ResponseWriter, r *http.Request)
}
