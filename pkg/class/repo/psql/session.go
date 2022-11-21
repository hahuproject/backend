package class_psql_repo

import class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"

func (sessionRepo PsqlClassRepoAdapter) StoreSession(session class_domain.Session) (class_domain.Session, error) {
	var addedSession class_domain.Session

	var course_user_id string
	sessionRepo.log.Println("Added successfully 0")
	sessionRepo.log.Println(session.Course.ID)
	sessionRepo.log.Println(session.Instructor.ID)
	err := sessionRepo.db.QueryRow("SELECT public.instructor_courses.course_user_id FROM public.instructor_courses WHERE course_id = $1 AND user_id = $2", session.Course.ID, session.Instructor.ID).Scan(&course_user_id)

	sessionRepo.log.Println(err)
	sessionRepo.log.Println(course_user_id)
	sessionRepo.log.Println("Added successfully 1")
	if err != nil || course_user_id == "" {
		sessionRepo.log.Println("Added successfully 2")
		return addedSession, err
	}
	sessionRepo.log.Println("Added successfully 3")

	err = sessionRepo.db.QueryRow("INSERT INTO public.sessions (duration, start_date, section_id, course_user_id) VALUES ($1,$2,$3,$4) RETURNING session_id, duration, start_date, section_id", session.Duration, session.StartDate, session.Section.ID, course_user_id).Scan(&addedSession.ID, &addedSession.Duration, &addedSession.StartDate, &addedSession.Section.ID)
	sessionRepo.log.Println("Added successfully 4")
	if err != nil {
		return addedSession, err
	}
	sessionRepo.log.Println("Added successfully 5")
	addedSession.Section, err = sessionRepo.FindSection(addedSession.Section.ID)
	sessionRepo.log.Println("Added successfully 6")

	return addedSession, err

}
func (sessionRepo PsqlClassRepoAdapter) FindSessions() ([]class_domain.Session, error) {

	var sessions []class_domain.Session = make([]class_domain.Session, 0)

	rows, err := sessionRepo.db.Query("SELECT * FROM public.sessions")
	if err != nil {
		return sessions, err
	}

	for rows.Next() {
		var session class_domain.Session
		rows.Scan(&session.ID, &session.Duration, &session.StartDate, &session.Section.ID)

		sessions = append(sessions, session)
	}

	return sessions, nil

	// for rows.Next() {
	// 	var session class_domain.Session
	// 	err = rows.Scan(&session.ID, &session.Duration, &session.StartDate, &session.Section.ID, &session.Section.Name, &session.Section.)
	// }

}
func (sessionRepo PsqlClassRepoAdapter) FindSessionBySection(sectionId string) ([]class_domain.Session, error) {
	var sessions []class_domain.Session = make([]class_domain.Session, 0)

	rows, err := sessionRepo.db.Query(`
	SELECT 
		public.sessions.session_id, public.sessions.duration, public.sessions.start_date, public.sessions.section_id,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.sessions 
		INNER JOIN public.instructor_courses ON public.instructor_courses.course_user_id = public.sessions.course_user_id
		INNER JOIN public.courses ON public.instructor_courses.course_id = public.courses.course_id
		INNER JOIN public.users ON public.instructor_courses.user_id = public.users.user_id
		INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
		INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE section_id = $1`, sectionId)
	if err != nil {
		return sessions, err
	}

	for rows.Next() {
		var session class_domain.Session
		rows.Scan(
			&session.ID, &session.Duration, &session.StartDate, &session.Section.ID,
			&session.Course.ID, &session.Course.Name, &session.Course.CreditHr, &session.Course.Color,
			&session.Instructor.ID, &session.Instructor.FirstName, &session.Instructor.LastName, &session.Instructor.Email, &session.Instructor.Phone, &session.Instructor.Username, &session.Instructor.ProfilePic, &session.Instructor.Verified, &session.Instructor.Type,
			&session.Instructor.Address.ID, &session.Instructor.Address.Country, &session.Instructor.Address.Region, &session.Instructor.Address.City, &session.Instructor.Address.SubCity, &session.Instructor.Address.Woreda, &session.Instructor.Address.HouseNo,
		)

		// session.Section, _ = sessionRepo.FindSection(session.Section.ID)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (sessionRepo PsqlClassRepoAdapter) FindSessionsByStudent(studentId string) ([]class_domain.Session, error) {
	var sessions []class_domain.Session = make([]class_domain.Session, 0)

	/*

		select  _session_id_
		_duration_
		_start_date_
		_course_id_
		_instructor_id_
		_section_id_

	*/

	rows, err := sessionRepo.db.Query(`
	SELECT
	public.sessions.session_id, public.sessions.duration, public.sessions.start_date,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no,
		public.sections.section_id
		FROM public.sessions
		INNER JOIN public.instructor_courses ON public.instructor_courses.course_user_id = public.sessions.course_user_id
		INNER JOIN public.courses ON public.instructor_courses.course_id = public.courses.course_id
		INNER JOIN public.users ON public.instructor_courses.user_id = public.users.user_id
		INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
		INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
		INNER JOIN public.sections ON public.sessions.section_id = public.sections.section_id
		INNER JOIN public.student_sections ON public.student_sections.section_id = public.sessions.section_id
		WHERE public.student_sections.student_id = $1`, studentId)
	if err != nil {
		return sessions, err
	}

	for rows.Next() {
		var session class_domain.Session
		rows.Scan(
			&session.ID, &session.Duration, &session.StartDate,
			&session.Course.ID, &session.Course.Name, &session.Course.CreditHr, &session.Course.Color,
			&session.Instructor.ID, &session.Instructor.FirstName, &session.Instructor.LastName, &session.Instructor.Email, &session.Instructor.Phone, &session.Instructor.Username, &session.Instructor.ProfilePic, &session.Instructor.Verified, &session.Instructor.Type,
			&session.Instructor.Address.ID, &session.Instructor.Address.Country, &session.Instructor.Address.Region, &session.Instructor.Address.City, &session.Instructor.Address.SubCity, &session.Instructor.Address.Woreda, &session.Instructor.Address.HouseNo,
			&session.Section.ID)

		session.Section, _ = sessionRepo.FindSection(session.Section.ID)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (sessionRepo PsqlClassRepoAdapter) FindSessionsByInstructor(instructorId string) ([]class_domain.Session, error) {
	var sessions []class_domain.Session = make([]class_domain.Session, 0)

	/*

		select  _session_id_
		_duration_
		_start_date_
		_course_id_
		_instructor_id_
		_section_id_

	*/

	rows, err := sessionRepo.db.Query(`
	SELECT
	public.sessions.session_id, public.sessions.duration, public.sessions.start_date,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no,
		public.sections.section_id
		FROM public.sessions 
		INNER JOIN public.instructor_courses ON public.instructor_courses.course_user_id = public.sessions.course_user_id
		INNER JOIN public.courses ON public.instructor_courses.course_id = public.courses.course_id
		INNER JOIN public.users ON public.instructor_courses.user_id = public.users.user_id
		INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
		INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
		INNER JOIN public.sections ON public.sessions.section_id = public.sections.section_id
		WHERE public.instructor_courses.user_id = $1`, instructorId)
	if err != nil {
		return sessions, err
	}

	for rows.Next() {
		var session class_domain.Session
		rows.Scan(
			&session.ID, &session.Duration, &session.StartDate,
			&session.Course.ID, &session.Course.Name, &session.Course.CreditHr, &session.Course.Color,
			&session.Instructor.ID, &session.Instructor.FirstName, &session.Instructor.LastName, &session.Instructor.Email, &session.Instructor.Phone, &session.Instructor.Username, &session.Instructor.ProfilePic, &session.Instructor.Verified, &session.Instructor.Type,
			&session.Instructor.Address.ID, &session.Instructor.Address.Country, &session.Instructor.Address.Region, &session.Instructor.Address.City, &session.Instructor.Address.SubCity, &session.Instructor.Address.Woreda, &session.Instructor.Address.HouseNo,
			&session.Section.ID)

		session.Section, _ = sessionRepo.FindSection(session.Section.ID)

		sessions = append(sessions, session)
	}

	sessionRepo.log.Println(len(sessions))

	return sessions, nil
}

func (repo PsqlClassRepoAdapter) DeleteSession(sessionId string) error {
	_, err := repo.db.Exec(`DELETE FROM public.sessions WHERE session_id = $1`, sessionId)
	return err
}
