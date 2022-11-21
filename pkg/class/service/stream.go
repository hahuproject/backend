package class_service

import (
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (service ClassServiceAdapter) AddStream(token string, stream class_domain.Stream) (class_domain.Stream, error) {
	var addedStream class_domain.Stream

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedStream, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return addedStream, class_utils.ErrUnauthorized
	}

	addedStream, err = service.repo.StoreStream(stream)

	return addedStream, err
}

func (service ClassServiceAdapter) UpdateStream(token string, stream class_domain.Stream) (class_domain.Stream, error) {
	var updatedStream class_domain.Stream

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedStream, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return updatedStream, class_utils.ErrUnauthorized
	}

	return service.repo.UpdateStream(stream)
}

func (service ClassServiceAdapter) DeleteStream(token, streamId string) error {

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return class_utils.ErrUnauthorized
	}

	err = service.repo.DeleteStream(streamId)

	return err
}
