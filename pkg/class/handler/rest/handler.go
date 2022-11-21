package class_rest_handler

import (
	"log"

	class_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/handler"
	class_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/service"
)

type RestClassHandlerAdapter struct {
	log          *log.Logger
	classService class_service.ClassServicePort
}

func NewRestClassHandlerAdapter(log *log.Logger, classService class_service.ClassServicePort) class_handler.ClassHandlerPort {
	return &RestClassHandlerAdapter{log: log, classService: classService}
}
