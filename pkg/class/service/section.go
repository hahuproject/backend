package class_service

import (
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (sectionService ClassServiceAdapter) AddSection(token string, section class_domain.Section) (class_domain.Section, error) {
	var addedSection class_domain.Section

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedSection, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return addedSection, class_utils.ErrUnauthorized
	}

	addedSection, err = sectionService.repo.StoreSection(section)

	return addedSection, err
}
func (sectionService ClassServiceAdapter) GetSections(token string) ([]class_domain.Section, error) {
	var sections []class_domain.Section
	sectionService.log.Println("GET SECTIONS SERVICE 00")

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		sectionService.log.Println("GET SECTIONS SERVICE 01")
		sectionService.log.Println(err)
		return sections, err
	}

	sectionService.log.Println("GET SECTIONS SERVICE 0")

	if user.Type == "SUPER_ADMIN" || user.Type == "ADMIN" || user.Type == "REGISTRY_OFFICER" {
		return sectionService.repo.FindSections()
	}

	sectionService.log.Println("GET SECTIONS SERVICE 1")
	if user.Type == "DEPARTMENT_HEAD" {
		return sectionService.repo.FindSectionsByDepartmentHead(user.ID)
	}

	sectionService.log.Println("GET SECTIONS SERVICE 2")
	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" {
		return sectionService.repo.FindSectionsByInstructor(user.ID)
	}

	sectionService.log.Println("GET SECTIONS SERVICE 3")
	return sections, nil
}
func (sectionService ClassServiceAdapter) GetSection(token, sectionId string) (class_domain.Section, error) {
	var section class_domain.Section

	_, err := class_utils.CheckAuth(token)
	if err != nil {
		return section, err
	}

	return sectionService.repo.FindSection(sectionId)
}

func (sectionService ClassServiceAdapter) AddStudentToSection(token, userId, sectionId string, courses []string) (class_domain.Section, error) {
	var section class_domain.Section

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return section, err
	}

	if user.Type != "DEPARTMENT_HEAD" && user.Type != "REGISTRY_OFFICER" && user.Type != "SUB_REGISTRY_OFFICER" {
		return section, class_utils.ErrUnauthorized
	}

	return sectionService.repo.AddStudentToSection(userId, sectionId, courses)
}
func (sectionService ClassServiceAdapter) RemoveStudentFromSection(token, userId, sectionId string, courses []string) (class_domain.Section, error) {
	var section class_domain.Section

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return section, err
	}

	if user.Type != "DEPARTMENT_HEAD" && user.Type != "REGISTRY_OFFICER" && user.Type != "SUB_REGISTRY_OFFICER" {
		return section, class_utils.ErrUnauthorized
	}

	return sectionService.repo.RemoveStudentFromSection(userId, sectionId, courses)
}

func (classService ClassServiceAdapter) UpdateSection(token string, section class_domain.Section) (class_domain.Section, error) {
	var updatedClass class_domain.Section

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedClass, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return updatedClass, class_utils.ErrUnauthorized
	}

	return classService.repo.UpdateSection(section)
}

func (classService ClassServiceAdapter) DeleteSection(token string, id string) error {

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return class_utils.ErrUnauthorized
	}

	return classService.repo.DeleteSection(id)
}
