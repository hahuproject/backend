package class_service

import (
	"errors"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

func (departmentService ClassServiceAdapter) AddDepartment(token string, department class_domain.Department) (class_domain.Department, error) {
	var addedDepartment class_domain.Department

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedDepartment, err
	}

	if user.Type != "SUPER_ADMIN" && user.Type != "ADMIN" {
		return addedDepartment, class_utils.ErrUnauthorized
	}

	if department.Name == "" {
		return addedDepartment, class_utils.ErrInvalidDepartmentName
	}

	addedDepartment, err = departmentService.repo.StoreDepartment(department)

	return addedDepartment, err
}
func (departmentService ClassServiceAdapter) GetDepartments(token string) ([]class_domain.Department, error) {
	var departments []class_domain.Department = make([]class_domain.Department, 0)
	var err error

	user, _ := class_utils.CheckAuth(token)

	// departmentService.log.Println("user.Type")
	// departmentService.log.Println(user.Type)

	// departmentService.log.Println(user.Type)

	if user.Type == "DEPARTMENT_HEAD" {
		departments, err = departmentService.repo.FindDepartmentsByHead(user)
	} else {
		// departmentService.log.Println("else - ")
		departments, err = departmentService.repo.FindDepartments()
	}

	// departmentService.log.Println(len(departments))

	return departments, err
}
func (departmentService ClassServiceAdapter) GetDepartment(token string, id string) (class_domain.Department, error) {
	var department class_domain.Department

	_, err := class_utils.CheckAuth(token)
	if err != nil {
		return department, err
	}

	department, err = departmentService.repo.FindDepartment(id)

	return department, err
}

func (departmentService ClassServiceAdapter) UpdateDepartment(token string, department class_domain.Department) (class_domain.Department, error) {

	var updatedDepartment class_domain.Department

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedDepartment, err
	}

	if user.Type != "ADMIN" {
		return updatedDepartment, class_utils.ErrUnauthorized
	}

	if department.ID == "" {
		return updatedDepartment, errors.New("provide a valid department")
	}

	return departmentService.repo.UpdateDepartment(department)
}

func (departmentService ClassServiceAdapter) DeleteDepartment(token, departmentId string) error {

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "ADMIN" {
		return class_utils.ErrUnauthorized
	}

	return departmentService.repo.DeleteDepartment(departmentId)
}
