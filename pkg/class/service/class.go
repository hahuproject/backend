package class_service

import (
	"log"
	"os"
	"strings"

	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
	class_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func (classService ClassServiceAdapter) AddClass(token string, class class_domain.Class) (class_domain.Class, error) {
	var addedClass class_domain.Class

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return addedClass, err
	}

	classService.log.Println(user.Type)

	if user.Type != "DEPARTMENT_HEAD" {
		return addedClass, class_utils.ErrUnauthorized
	}

	classService.log.Println(user.ID == class.Department.Head.ID)
	classService.log.Println(user.ID)
	classService.log.Println(class.Department.Head.ID)

	if user.ID != class.Department.Head.ID {
		return addedClass, class_utils.ErrUnauthorized
	}

	addedClass, err = classService.repo.StoreClass(class)

	return addedClass, err
}

func (classService ClassServiceAdapter) GetClasses(token string) ([]class_domain.Class, error) {
	var classes []class_domain.Class

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return classes, err
	}

	if user.Type == "SUPER_ADMIN" || user.Type == "ADMIN" || user.Type == "REGISTRY_OFFICER" {
		classes, err = classService.repo.FindClasses()
	}

	if user.Type == "DEPARTMENT_HEAD" {
		classes, err = classService.repo.FindClassesByDepartmentHead(user.ID)
	}

	return classes, err
}

func (classService ClassServiceAdapter) GetClass(token string, id string) (class_domain.Class, error) {
	var class class_domain.Class

	_, err := class_utils.CheckAuth(token)
	if err != nil {
		return class, err
	}

	return classService.repo.FindClass(id)
}

func (classService ClassServiceAdapter) UpdateClass(token string, class class_domain.Class) (class_domain.Class, error) {
	var updatedClass class_domain.Class

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedClass, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return updatedClass, class_utils.ErrUnauthorized
	}

	return classService.repo.UpdateClass(class)
}

func (classService ClassServiceAdapter) DeleteClass(token string, id string) error {

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return class_utils.ErrUnauthorized
	}

	return classService.repo.DeleteClass(id)
}

func (classService ClassServiceAdapter) AddCourseToClass(token, classId, courseId string) (class_domain.Class, error) {
	var updatedClass class_domain.Class

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedClass, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return updatedClass, class_utils.ErrUnauthorized
	}

	return classService.repo.StoreClassCourse(classId, courseId)
}
func (classService ClassServiceAdapter) RemoveCourseFromClass(token, classId, courseId string) (class_domain.Class, error) {
	var updatedClass class_domain.Class

	user, err := class_utils.CheckAuth(token)
	if err != nil {
		return updatedClass, err
	}

	if user.Type != "DEPARTMENT_HEAD" {
		return updatedClass, class_utils.ErrUnauthorized
	}

	return classService.repo.RemoveClassCourse(classId, courseId)
}

func (classService ClassServiceAdapter) GetClassMasterSheet(token, classId string) (string, error) {

	// user, err := class_utils.CheckAuth(token)
	// if err != nil {
	// 	return "", class_utils.ErrUnauthorized
	// }

	// if user.Type != "DEPARTMENT_HEAD" {
	// 	return "", class_utils.ErrUnauthorized
	// }

	/*

		Data to generate master sheet

		- Section
			- class - name
			- year

	*/

	//Generate PDF

	m := pdf.NewMaroto(consts.Landscape, consts.A4)

	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("MASTER SHEET FOR MIDDLE LEVEL IV TVET TRAINING PROGRAMME 2011 E.C", props.Text{
				Size:  8,
				Style: consts.Bold,
				Align: consts.Center,
				Top:   0,
			})
		})
	})

	m.Row(8, func() {
		m.Col(12, func() {
			m.Text("ADDIS ABABA TEGBAR - ID TEVT COLLEGE", props.Text{
				Size:  8,
				Style: consts.Bold,
				Align: consts.Top,
				Top:   0,
			})
		})
	})

	m.Row(12, func() {
		m.Col(8, func() {
			m.Text("Occupational Title :  Mechatronics and Instrumentation Servicing Management", props.Text{
				Size:  8,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   4,
			})
		})
		m.Col(2, func() {
			m.Text("Program :Extension", props.Text{
				Size:  8,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   4,
			})
		})
		m.Col(2, func() {
			m.Text("Level : IV", props.Text{
				Size:  8,
				Style: consts.Normal,
				Align: consts.Left,
				Top:   4,
			})
		})
	})

	_, err := os.Open("/uploads/master-sheets")
	if err != nil {
		err = os.MkdirAll("uploads/master-sheets", 0755)
		if err != nil {
			log.Println(err)
		}
	}

	m.OutputFileAndClose("./uploads/master-sheets/" + strings.ReplaceAll(classId, "/", "-") + ".pdf")

	return "/uploads/master-sheets/" + strings.ReplaceAll(classId, "/", "-") + ".pdf", nil

	// return classService.repo.RemoveClassCourse(classId, courseId)
	// return " ", nil
}
