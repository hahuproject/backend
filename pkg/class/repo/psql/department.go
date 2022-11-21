package class_psql_repo

import (
	"database/sql"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

func (departmentRepo PsqlClassRepoAdapter) StoreDepartment(department class_domain.Department) (class_domain.Department, error) {
	var addedDepartment class_domain.Department
	err := departmentRepo.db.QueryRow("INSERT INTO public.departments (name) VALUES ($1) RETURNING department_id, name", department.Name).Scan(&addedDepartment.ID, &addedDepartment.Name)
	return addedDepartment, err
}

func queryDepartmentHead(db *sql.DB, id string) (auth_domain.DepartmentHead, error) {
	var _user auth_domain.DepartmentHead
	err := db.QueryRow("SELECT public.users.user_id, public.department_heads.user_id, first_name, last_name, email, phone, username, verified, type FROM public.users INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id INNER JOIN public.department_heads ON public.users.user_id = public.department_heads.user_id WHERE public.users.user_id = $1", id).Scan(&_user.ID, &_user.User.ID, &_user.User.FirstName, &_user.User.LastName, &_user.User.Email, &_user.User.Phone, &_user.User.Username, &_user.User.Verified, &_user.User.Type)
	if err != nil {
		return _user, err
	}

	return _user, nil
}

func populateDepartment(db *sql.DB, department *class_domain.Department) {
	department.Classes = make([]class_domain.Class, 0)
	department.Streams = make([]class_domain.Stream, 0)

	classRows, err := db.Query("SELECT class_id, name FROM public.classes WHERE department_id = $1", department.ID)
	if err == nil {
		for classRows.Next() {
			var class class_domain.Class
			classRows.Scan(&class.ID, &class.Name)

			department.Classes = append(department.Classes, class)
		}

	}

	streamRows, err := db.Query(`SELECT stream_id, name FROM public.streams WHERE department_id = $1`, department.ID)
	if err == nil {
		for streamRows.Next() {
			var stream class_domain.Stream
			streamRows.Scan(&stream.ID, &stream.Name)
			stream.Department = *department

			department.Streams = append(department.Streams, stream)
		}

	}
}

func (departmentRepo PsqlClassRepoAdapter) FindDepartments() ([]class_domain.Department, error) {
	var departments []class_domain.Department = make([]class_domain.Department, 0)

	rows, err := departmentRepo.db.Query("SELECT department_id, name, head FROM public.departments ")
	if err != nil {
		return departments, err
	}

	defer rows.Close()

	for rows.Next() {
		var department class_domain.Department
		var nullHead sql.NullString
		rows.Scan(&department.ID, &department.Name, &nullHead)
		if nullHead.Valid {
			department.Head, _ = queryDepartmentHead(departmentRepo.db, nullHead.String)
		}

		populateDepartment(departmentRepo.db, &department)

		departments = append(departments, department)
	}
	// departmentRepo.log.Println(len(departments))

	return departments, nil
}

func (departmentRepo PsqlClassRepoAdapter) FindDepartmentsByHead(head auth_domain.User) ([]class_domain.Department, error) {
	var departments []class_domain.Department = make([]class_domain.Department, 0)

	rows, err := departmentRepo.db.Query("SELECT department_id, name, head FROM public.departments WHERE head = $1", head.ID)
	if err != nil {
		return departments, err
	}

	defer rows.Close()

	for rows.Next() {
		var department class_domain.Department
		var nullHead sql.NullString
		rows.Scan(&department.ID, &department.Name, &nullHead)
		if nullHead.Valid {
			department.Head, _ = queryDepartmentHead(departmentRepo.db, nullHead.String)
		}

		department.Classes = make([]class_domain.Class, 0)

		populateDepartment(departmentRepo.db, &department)

		departments = append(departments, department)
	}

	return departments, nil
}

func (departmentRepo PsqlClassRepoAdapter) FindDepartment(id string) (class_domain.Department, error) {
	var department class_domain.Department
	var nullHead sql.NullString
	err := departmentRepo.db.QueryRow("SELECT * FROM public.departments WHERE department_id = $1", id).Scan(&department.ID, &department.Name, &nullHead)
	if nullHead.Valid {
		department.Head, _ = queryDepartmentHead(departmentRepo.db, nullHead.String)
	}
	return department, err
}

func (departmentRepo PsqlClassRepoAdapter) UpdateDepartment(department class_domain.Department) (class_domain.Department, error) {
	var updatedDepartment class_domain.Department

	var departmentHeadId sql.NullString
	if department.Head.ID != "" {
		departmentHeadId.Valid = true
		departmentHeadId.String = department.Head.ID
	}

	err := departmentRepo.db.QueryRow("UPDATE public.departments SET name = $1, head = $2 WHERE department_id = $3 RETURNING department_id, name, head", department.Name, departmentHeadId, department.ID).Scan(&updatedDepartment.ID, &updatedDepartment.Name, &departmentHeadId)
	if err != nil {
		return updatedDepartment, err
	}

	return departmentRepo.FindDepartment(updatedDepartment.ID)
}

func (departmentRepo PsqlClassRepoAdapter) DeleteDepartment(departmentId string) error {

	_, err := departmentRepo.db.Query("DELETE FROM public.departments WHERE department_id = $1", departmentId)
	if err != nil {
		return err
	}

	return nil
}
