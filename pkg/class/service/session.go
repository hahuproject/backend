package class_service

import (
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (sessionService ClassServiceAdapter) AddSession(token string, session class_domain.Session) (class_domain.Session, error) {
	var addedSession class_domain.Session

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedSession, err
	}

	if user.Type != "DEPARTMENT_HEAD" && user.Type != "INSTRUCTOR" && user.Type != "SUPERVISOR" {
		return addedSession, class_utils.ErrUnauthorized
	}

	return sessionService.repo.StoreSession(session)
}

func (sessionService ClassServiceAdapter) GetSessions(token string) ([]class_domain.Session, error) {
	var studentSessions []class_domain.Session = make([]class_domain.Session, 0)

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return studentSessions, err
	}

	if user.Type == "ADMIN" || user.Type == "REGISTRY_OFFICER" || user.Type == "DEPARTMENT_HEAD" {
		return sessionService.repo.FindSessions()
	}

	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" {
		return sessionService.repo.FindSessionsByInstructor(user.ID)

	}

	if user.Type == "STUDENT" {
		return sessionService.repo.FindSessionsByStudent(user.ID)
	}

	return studentSessions, class_utils.ErrUnauthorized
}

func (service ClassServiceAdapter) DeleteSession(token, sessionId string) error {
	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return class_utils.ErrUnauthorized
	}

	return service.repo.DeleteSession(sessionId)
}
