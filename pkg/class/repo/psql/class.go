package class_psql_repo

import (
	"database/sql"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

func populateClass(db *sql.DB, class *class_domain.Class) {
	//course , section, students

	//course
	class.Courses = make([]class_domain.Course, 0)
	rows, err := db.Query(`SELECT
	public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color
	FROM public.class_courses
	INNER JOIN public.courses ON public.class_courses.course_id = public.courses.course_id
	WHERE public.class_courses.class_id = $1`, class.ID)

	if err == nil {
		for rows.Next() {
			var course class_domain.Course
			rows.Scan(&course.ID, &course.Name, &course.CreditHr, &course.Color)

			populateCourse(db, &course)

			class.Courses = append(class.Courses, course)
		}

		defer rows.Close()
	}

	//Sections
	class.Sections = make([]class_domain.Section, 0)
	rows, err = db.Query(`SELECT
	public.sections.section_id, public.sections.name
	FROM public.sections
	WHERE public.sections.class_id = $1`, class.ID)

	if err == nil {
		for rows.Next() {
			var section class_domain.Section
			rows.Scan(&section.ID, &section.Name)

			class.Sections = append(class.Sections, section)
		}

		defer rows.Close()
	}

	class.Students = make([]auth_domain.Student, 0)

	// Students
	rows, err = db.Query(`SELECT
	public.students.user_id, public.students.middle_name, public.students.gender, public.students.disability, 
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.student_sections
	INNER JOIN public.sections ON public.sections.section_id = public.student_sections.section_id
	INNER JOIN public.students ON public.students.user_id = public.student_sections.student_id
	INNER JOIN public.users ON public.students.user_id = public.users.user_id
	INNER JOIN public.user_types 
	ON public.users.user_id = public.user_types.user_id 
	INNER JOIN public.addresses
	ON public.users.address = public.addresses.address_id
	WHERE public.sections.class_id = $1 AND public.users.verified = $2`, class.ID, true)

	if err == nil {
		for rows.Next() {
			var student auth_domain.Student
			rows.Scan(
				&student.ID, &student.MiddleName, &student.Gender, &student.Disablility,
				&student.User.ID, &student.User.FirstName, &student.User.LastName, &student.User.Email, &student.User.Phone, &student.User.Username, &student.User.ProfilePic, &student.User.Verified, &student.User.Type,
				&student.User.Address.ID, &student.User.Address.Country, &student.User.Address.Region, &student.User.Address.City, &student.User.Address.SubCity, &student.User.Address.Woreda, &student.User.Address.HouseNo,
			)

			class.Students = append(class.Students, student)
		}

		defer rows.Close()
	}

}

func (classRepo PsqlClassRepoAdapter) StoreClass(class class_domain.Class) (class_domain.Class, error) {
	var addedClass class_domain.Class

	// classRepo.log.Println(class.SubName)

	err := classRepo.db.QueryRow("INSERT INTO public.classes (name, sub_name, department_id, stream_id) VALUES ($1,$2, $3, $4) RETURNING class_id", class.Name, sql.NullString{Valid: class.SubName != "", String: class.SubName}, class.Department.ID, sql.NullString{Valid: class.Stream.ID != "", String: class.Stream.ID}).Scan(&addedClass.ID)
	if err != nil {
		classRepo.log.Println("err in insert", err)
		return addedClass, err
	}
	addedClass, err = classRepo.FindClass(addedClass.ID)

	return addedClass, err
}

func (classRepo PsqlClassRepoAdapter) FindClasses() ([]class_domain.Class, error) {
	var classes []class_domain.Class = make([]class_domain.Class, 0)

	rows, err := classRepo.db.Query(`
	SELECT 
	public.classes.class_id, public.classes.name,public.classes.sub_name,
	public.departments.department_id, public.departments.name,
	public.department_heads.user_id,
	public.streams.stream_id, public.streams.name,
	public.users.user_id, public.users.first_name, public.users.last_name, public.users.email, public.users.phone, public.users.profile_pic, public.users.username, public.users.verified,
	public.addresses.address_id, public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity, public.addresses.woreda, public.addresses.house_no
	FROM public.classes 
	INNER JOIN public.departments ON public.classes.department_id = public.departments.department_id 
	LEFT JOIN public.streams ON public.classes.stream_id = public.streams.stream_id 
	LEFT JOIN public.department_heads ON public.departments.head = public.department_heads.user_id
	LEFT JOIN public.users ON public.department_heads.user_id = public.users.user_id
	LEFT JOIN public.addresses ON public.users.address = public.addresses.address_id`)
	if err != nil {
		return classes, err
	}

	defer rows.Close()

	for rows.Next() {
		var class class_domain.Class
		var _subName sql.NullString
		rows.Scan(
			&class.ID, &class.Name, &_subName,
			&class.Department.ID, &class.Department.Name,
			&class.Department.Head.ID,
			&class.Stream.ID, &class.Stream.Name,
			&class.Department.Head.User.ID, &class.Department.Head.User.FirstName, &class.Department.Head.User.LastName, &class.Department.Head.User.Email, &class.Department.Head.User.Phone, &class.Department.Head.User.ProfilePic, &class.Department.Head.User.Username, &class.Department.Head.User.Verified,
			&class.Department.Head.User.Address.ID, &class.Department.Head.User.Address.Country, &class.Department.Head.User.Address.Region, &class.Department.Head.User.Address.City, &class.Department.Head.User.Address.SubCity, &class.Department.Head.User.Address.Woreda, &class.Department.Head.User.Address.HouseNo)

		if _subName.Valid {
			class.SubName = _subName.String
		}

		populateClass(classRepo.db, &class)

		classes = append(classes, class)
	}

	return classes, nil
}
func (classRepo PsqlClassRepoAdapter) FindClass(id string) (class_domain.Class, error) {
	var class class_domain.Class
	var _streamIdNull sql.NullString
	var _streamNameNull sql.NullString
	var _subName sql.NullString
	err := classRepo.db.QueryRow(`SELECT 
	public.classes.class_id, public.classes.name,public.classes.sub_name, 
	public.departments.department_id, public.departments.name,
	public.department_heads.user_id,
	public.streams.stream_id, public.streams.name,
	public.users.user_id, public.users.first_name, public.users.last_name, public.users.email, public.users.phone, public.users.profile_pic, public.users.username, public.users.verified,
	public.addresses.address_id, public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity, public.addresses.woreda, public.addresses.house_no
	FROM public.classes 
	INNER JOIN public.departments ON public.classes.department_id = public.departments.department_id 
	LEFT JOIN public.streams ON public.classes.stream_id = public.streams.stream_id 
	LEFT JOIN public.department_heads ON public.departments.head = public.department_heads.user_id
	LEFT JOIN public.users ON public.department_heads.user_id = public.users.user_id
	LEFT JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.classes.class_id = $1`, id).Scan(
		&class.ID, &class.Name, &_subName,
		&class.Department.ID, &class.Department.Name,
		&class.Department.Head.ID,
		&_streamIdNull, &_streamNameNull,
		&class.Department.Head.User.ID, &class.Department.Head.User.FirstName, &class.Department.Head.User.LastName, &class.Department.Head.User.Email, &class.Department.Head.User.Phone, &class.Department.Head.User.ProfilePic, &class.Department.Head.User.Username, &class.Department.Head.User.Verified,
		&class.Department.Head.User.Address.ID, &class.Department.Head.User.Address.Country, &class.Department.Head.User.Address.Region, &class.Department.Head.User.Address.City, &class.Department.Head.User.Address.SubCity, &class.Department.Head.User.Address.Woreda, &class.Department.Head.User.Address.HouseNo)

	if _streamIdNull.Valid {
		class.Stream.ID = _streamIdNull.String
	}

	if _streamNameNull.Valid {
		class.Stream.Name = _streamNameNull.String
	}

	if _subName.Valid {
		class.SubName = _subName.String
	}

	populateClass(classRepo.db, &class)

	// classRepo.log.Println("err in find", err)

	return class, err

}

func (classRepo PsqlClassRepoAdapter) FindClassesByDepartmentHead(departmentHeadId string) ([]class_domain.Class, error) {
	var classes []class_domain.Class = make([]class_domain.Class, 0)

	rows, err := classRepo.db.Query(`
	SELECT 
	public.classes.class_id, public.classes.name,public.classes.sub_name,
	public.departments.department_id, public.departments.name,
	public.department_heads.user_id,
	public.streams.stream_id, public.streams.name,
	public.users.user_id, public.users.first_name, public.users.last_name, public.users.email, public.users.phone, public.users.profile_pic, public.users.username, public.users.verified,
	public.addresses.address_id, public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity, public.addresses.woreda, public.addresses.house_no
	FROM public.classes 
	INNER JOIN public.departments ON public.classes.department_id = public.departments.department_id 
	LEFT JOIN public.streams ON public.classes.stream_id = public.streams.stream_id 
	INNER JOIN public.department_heads ON public.departments.head = public.department_heads.user_id
	INNER JOIN public.users ON public.department_heads.user_id = public.users.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.department_heads.user_id = $1`, departmentHeadId)
	if err != nil {
		return classes, err
	}

	defer rows.Close()

	for rows.Next() {
		var class class_domain.Class
		var _streamIdNull sql.NullString
		var _streamNameNull sql.NullString
		var _subName sql.NullString
		rows.Scan(
			&class.ID, &class.Name, &_subName,
			&class.Department.ID, &class.Department.Name,
			&class.Department.Head.ID,
			&_streamIdNull, &_streamNameNull,
			&class.Department.Head.User.ID, &class.Department.Head.User.FirstName, &class.Department.Head.User.LastName, &class.Department.Head.User.Email, &class.Department.Head.User.Phone, &class.Department.Head.User.ProfilePic, &class.Department.Head.User.Username, &class.Department.Head.User.Verified,
			&class.Department.Head.User.Address.ID, &class.Department.Head.User.Address.Country, &class.Department.Head.User.Address.Region, &class.Department.Head.User.Address.City, &class.Department.Head.User.Address.SubCity, &class.Department.Head.User.Address.Woreda, &class.Department.Head.User.Address.HouseNo)

		if _streamIdNull.Valid {
			class.Stream.ID = _streamIdNull.String
		}

		if _streamNameNull.Valid {
			class.Stream.Name = _streamNameNull.String
		}

		if _subName.Valid {
			class.SubName = _subName.String
		}

		populateClass(classRepo.db, &class)

		classes = append(classes, class)
	}

	return classes, nil
}
func (classRepo PsqlClassRepoAdapter) UpdateClass(class class_domain.Class) (class_domain.Class, error) {
	var updatedClass class_domain.Class
	err := classRepo.db.QueryRow("UPDATE public.classes SET name = $1, sub_name = $2 WHERE class_id = $3 RETURNING public.classes.class_id ", class.Name, sql.NullString{Valid: class.SubName != "", String: class.SubName}, class.ID).Scan(&updatedClass.ID)
	if err != nil {
		return updatedClass, err
	}
	updatedClass, err = classRepo.FindClass(updatedClass.ID)
	return updatedClass, err
}
func (classRepo PsqlClassRepoAdapter) DeleteClass(id string) error {
	_, err := classRepo.db.Query("DELETE FROM public.classes WHERE class_id = $1", id)
	return err
}

func (classRepo PsqlClassRepoAdapter) StoreClassCourse(classId, courseId string) (class_domain.Class, error) {
	var updatedClass class_domain.Class

	_, err := classRepo.db.Query("INSERT INTO public.class_courses (class_id, course_id) VALUES ($1, $2)", classId, courseId)
	if err != nil {
		return updatedClass, err
	}

	updatedClass, _ = classRepo.FindClass(classId)

	return updatedClass, nil
}
func (classRepo PsqlClassRepoAdapter) RemoveClassCourse(classId, courseId string) (class_domain.Class, error) {
	var updatedClass class_domain.Class

	_, err := classRepo.db.Query("DELETE FROM public.class_courses WHERE class_id = $1 AND course_id = $2", classId, courseId)
	if err != nil {
		return updatedClass, err
	}

	updatedClass, _ = classRepo.FindClass(classId)

	return updatedClass, nil
}
