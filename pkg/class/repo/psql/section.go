package class_psql_repo

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

func populateSectionByInstructor(instructorId string, sectionRepo PsqlClassRepoAdapter, section *class_domain.Section) {
	section.Class, _ = sectionRepo.FindClass(section.Class.ID)
	section.Sessions = make([]class_domain.Session, 0)
	section.Sessions, _ = sectionRepo.FindSessionsByInstructor(instructorId)
	section.Students = make([]auth_domain.Student, 0)
	rows, err := sectionRepo.db.Query(`
	SELECT
		public.students.user_id, public.students.middle_name, public.students.gender, public.students.disability, 
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.student_sections
		INNER JOIN public.students ON public.students.user_id = public.student_sections.student_id
		INNER JOIN public.users ON public.students.user_id = public.users.user_id
		INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
		INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.student_sections.section_id = $1`, section.ID)
	if err == nil {

		for rows.Next() {
			var student auth_domain.Student
			rows.Scan(&student.ID, &student.MiddleName, &student.Gender, &student.Disablility,
				&student.User.ID, &student.User.FirstName, &student.User.LastName, &student.User.Email, &student.User.Phone, &student.User.Username, &student.User.ProfilePic, &student.User.Verified, &student.User.Type,
				&student.User.Address.ID, &student.User.Address.Country, &student.User.Address.Region, &student.User.Address.City, &student.User.Address.SubCity, &student.User.Address.Woreda, &student.User.Address.HouseNo)

			section.Students = append(section.Students, student)
		}

		defer rows.Close()
	}
}

func populateSection(sectionRepo PsqlClassRepoAdapter, section *class_domain.Section) {
	section.Class, _ = sectionRepo.FindClass(section.Class.ID)
	section.Sessions = make([]class_domain.Session, 0)
	section.Sessions, _ = sectionRepo.FindSessionBySection(section.ID)
	section.Students = make([]auth_domain.Student, 0)
	sectionRepo.log.Println("Populate section")
	rows, err := sectionRepo.db.Query(`
	SELECT
		public.students.user_id, public.students.middle_name, public.students.gender, public.students.disability, 
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.student_sections
		INNER JOIN public.students ON public.students.user_id = public.student_sections.student_id
		INNER JOIN public.users ON public.students.user_id = public.users.user_id
		INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
		INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.student_sections.section_id = $1`, section.ID)
	if err != nil {
		sectionRepo.log.Println(err)
	}
	if err == nil {
		sectionRepo.log.Println("Populate section 1")
		for rows.Next() {
			sectionRepo.log.Println("Populate section 2")
			var student auth_domain.Student
			rows.Scan(&student.ID, &student.MiddleName, &student.Gender, &student.Disablility,
				&student.User.ID, &student.User.FirstName, &student.User.LastName, &student.User.Email, &student.User.Phone, &student.User.Username, &student.User.ProfilePic, &student.User.Verified, &student.User.Type,
				&student.User.Address.ID, &student.User.Address.Country, &student.User.Address.Region, &student.User.Address.City, &student.User.Address.SubCity, &student.User.Address.Woreda, &student.User.Address.HouseNo)

			section.Students = append(section.Students, student)
		}

		sectionRepo.log.Println(len(section.Students))
		defer rows.Close()
	}
}

func (sectionRepo PsqlClassRepoAdapter) StoreSection(section class_domain.Section) (class_domain.Section, error) {
	var addedSection class_domain.Section
	err := sectionRepo.db.QueryRow("INSERT INTO public.sections (name, year, class_id) VALUES ($1,$2, $3) RETURNING section_id, name,year, class_id", section.Name, section.Year, section.Class.ID).Scan(&addedSection.ID, &addedSection.Name, &addedSection.Year, &addedSection.Class.ID)
	if err != nil {
		return addedSection, err
	}

	populateSection(sectionRepo, &addedSection)

	return addedSection, err
}

func (sectionRepo PsqlClassRepoAdapter) FindSections() ([]class_domain.Section, error) {
	var sections []class_domain.Section = make([]class_domain.Section, 0)

	rows, err := sectionRepo.db.Query("SELECT section_id, name, year, class_id FROM public.sections")
	if err != nil {
		return sections, err
	}

	defer rows.Close()

	for rows.Next() {
		var section class_domain.Section
		rows.Scan(&section.ID, &section.Name, &section.Year, &section.Class.ID)

		populateSection(sectionRepo, &section)

		sections = append(sections, section)
	}

	return sections, nil
}

func (sectionRepo PsqlClassRepoAdapter) FindSectionsByClass(classId string) ([]class_domain.Section, error) {
	var sections []class_domain.Section = make([]class_domain.Section, 0)
	rows, err := sectionRepo.db.Query("SELECT section_id, name, year, class_id FROM public.sections WHERE class_id = $1", classId)
	if err != nil {
		return sections, err
	}

	for rows.Next() {
		var section class_domain.Section
		rows.Scan(&section.ID, &section.Name, &section.Year, &section.Class.ID)

		populateSection(sectionRepo, &section)

		sections = append(sections, section)
	}

	return sections, err
}

func (sectionRepo PsqlClassRepoAdapter) FindSectionsByDepartmentHead(departmentHeadId string) ([]class_domain.Section, error) {
	var sections []class_domain.Section = make([]class_domain.Section, 0)
	rows, err := sectionRepo.db.Query("SELECT public.sections.section_id, public.sections.name, public.sections.year, public.sections.class_id FROM public.sections INNER JOIN public.classes ON public.classes.class_id = public.sections.class_id INNER JOIN public.departments ON public.departments.department_id = public.classes.department_id INNER JOIN public.users ON public.users.user_id = public.departments.head WHERE public.users.user_id = $1", departmentHeadId)

	for rows.Next() {
		var section class_domain.Section
		rows.Scan(&section.ID, &section.Name, &section.Year, &section.Class.ID)

		populateSection(sectionRepo, &section)

		sections = append(sections, section)
	}

	return sections, err
}

func (sectionRepo PsqlClassRepoAdapter) FindSectionsByInstructor(instructorId string) ([]class_domain.Section, error) {
	var sections []class_domain.Section = make([]class_domain.Section, 0)
	rows, err := sectionRepo.db.Query(`
	SELECT 
		public.sections.section_id, public.sections.name, public.sections.year, public.sections.class_id 
	FROM public.sections 
		INNER JOIN public.sessions ON public.sessions.section_id = public.sections.section_id
		INNER JOIN public.instructor_courses ON public.instructor_courses.course_user_id = public.sessions.course_user_id
	WHERE public.instructor_courses.user_id = $1
	GROUP BY public.sections.section_id`, instructorId)

	for rows.Next() {
		var section class_domain.Section
		rows.Scan(&section.ID, &section.Name, &section.Year, &section.Class.ID)

		// populateSection(sectionRepo, &section)
		populateSectionByInstructor(instructorId, sectionRepo, &section)
		populateClass(sectionRepo.db, &section.Class)

		sections = append(sections, section)
	}

	return sections, err
}

func (sectionRepo PsqlClassRepoAdapter) FindSection(id string) (class_domain.Section, error) {
	var section class_domain.Section
	err := sectionRepo.db.QueryRow("SELECT section_id, name, year, class_id FROM public.sections WHERE section_id = $1", id).Scan(&section.ID, &section.Name, &section.Year, &section.Class.ID)
	populateSection(sectionRepo, &section)
	return section, err
}

func (sectionRepo PsqlClassRepoAdapter) AddStudentToSection(userId, sectionId string, courses []string) (class_domain.Section, error) {
	var section class_domain.Section

	err := sectionRepo.db.QueryRow(`
	INSERT INTO public.student_sections 
	(student_id, section_id) 
	VALUES ($1,$2) RETURNING section_id`, userId, sectionId).Scan(&section.ID)
	if err != nil {
		return section, err
	}

	for i := 0; i < len(courses); i++ {
		_, err = sectionRepo.db.Query(`INSERT INTO public.student_courses (student_id, course_id) VALUES ($1,$2)`, userId, courses[i])
	}

	return sectionRepo.FindSection(sectionId)
}
func (sectionRepo PsqlClassRepoAdapter) RemoveStudentFromSection(userId, sectionId string, courses []string) (class_domain.Section, error) {
	var section class_domain.Section

	// sectionRepo.log.Println("userId")
	// sectionRepo.log.Println(userId)
	// sectionRepo.log.Println(sectionId)

	_, err := sectionRepo.db.Query(`DELETE FROM public.student_sections WHERE user_id = $1 AND section_id = $2`, userId, sectionId)
	if err != nil {
		return section, err
	}

	for i := 0; i < len(courses); i++ {
		_, err = sectionRepo.db.Query(`DELETE FROM public.student_courses WHERE student_id = $1 AND course_id != $2`, userId, courses[i])
	}

	return sectionRepo.FindSection(sectionId)
}

func (classRepo PsqlClassRepoAdapter) UpdateSection(section class_domain.Section) (class_domain.Section, error) {
	var updatedSection class_domain.Section
	err := classRepo.db.QueryRow("UPDATE public.sections SET name = $1, year = $2 WHERE section_id = $3 RETURNING section_id", section.Name, section.Year, section.ID).Scan(&updatedSection.ID)
	updatedSection, _ = classRepo.FindSection(updatedSection.ID)
	return updatedSection, err
}
func (classRepo PsqlClassRepoAdapter) DeleteSection(id string) error {
	_, err := classRepo.db.Query("DELETE FROM public.sections WHERE section_id = $1", id)
	return err
}
