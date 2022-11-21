package class_psql_repo

import (
	"database/sql"
	"log"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

func populateCourse(db *sql.DB, course *class_domain.Course) {

	course.Users = make([]auth_domain.User, 0)
	course.Prerequisites = make([]class_domain.Course, 0)

	rows, err := db.Query(`SELECT
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.instructor_courses
	INNER JOIN public.users ON public.instructor_courses.user_id = public.users.user_id
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id 
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.instructor_courses.course_id = $1`, course.ID)

	if err != nil {
		log.Println("populate course")
		log.Println(err)
	}
	if err == nil {
		for rows.Next() {
			var user auth_domain.User
			rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.ProfilePic, &user.Verified, &user.Type,
				&user.Address.ID, &user.Address.Country, &user.Address.Region, &user.Address.City, &user.Address.SubCity, &user.Address.Woreda, &user.Address.HouseNo)

			course.Users = append(course.Users, user)
		}

		defer rows.Close()
	}

	coursePrerequisites, err := db.Query(`
	SELECT public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color
	FROM public.course_prerequisites
	INNER JOIN public.courses ON prerequisite = public.courses.course_id
	WHERE public.course_prerequisites.course_id = $1`, course.ID)

	if err == nil {
		for coursePrerequisites.Next() {
			var _course class_domain.Course
			coursePrerequisites.Scan(&_course.ID, &_course.Name, &_course.CreditHr, &_course.Color)

			course.Prerequisites = append(course.Prerequisites, _course)
		}

		defer coursePrerequisites.Close()
	}

	// log.Println(course.Prerequisites)
}

func (courseRepo PsqlClassRepoAdapter) StoreCourse(course class_domain.Course) (class_domain.Course, error) {
	var addedCourse class_domain.Course
	err := courseRepo.db.QueryRow("INSERT INTO public.courses (name, credit_hr, color) VALUES ($1,$2, $3) RETURNING *", course.Name, course.CreditHr, course.Color).Scan(&addedCourse.ID, &addedCourse.Name, &addedCourse.CreditHr, &addedCourse.Color)

	if err != nil {
		return addedCourse, err
	}

	for i := 0; i < len(course.Prerequisites); i++ {
		courseRepo.db.Query("INSERT INTO public.course_prerequisites (course_id, prerequisite) VALUES ($1,$2) RETURNING *", addedCourse.ID, course.Prerequisites[i].ID)
	}

	populateCourse(courseRepo.db, &addedCourse)

	return addedCourse, err
}

func (courseRepo PsqlClassRepoAdapter) FindCourses() ([]class_domain.Course, error) {
	var courses []class_domain.Course = make([]class_domain.Course, 0)
	rows, err := courseRepo.db.Query("SELECT * FROM public.courses")
	if err != nil {
		return courses, err
	}

	for rows.Next() {
		var course class_domain.Course
		_ = rows.Scan(&course.ID, &course.Name, &course.CreditHr, &course.Color)

		populateCourse(courseRepo.db, &course)

		courses = append(courses, course)
	}

	return courses, nil
}

func (courseRepo PsqlClassRepoAdapter) FindCourse(id string) (class_domain.Course, error) {
	var course class_domain.Course
	err := courseRepo.db.QueryRow("SELECT * FROM public.courses WHERE course_id = $1", id).Scan(&course.ID, &course.Name, &course.CreditHr, &course.Color)
	if err != nil {
		return course, err
	}

	populateCourse(courseRepo.db, &course)
	return course, nil
}

func (courseRepo PsqlClassRepoAdapter) FindCoursesByClass(classId string) ([]class_domain.Course, error) {
	var courses []class_domain.Course
	rows, err := courseRepo.db.Query("SELECT public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color FROM public.class_courses INNER JOIN public.courses ON public.courses.course_id = public.class_courses.course_id WHERE class_id = $1", classId)
	if err != nil {
		return courses, err
	}

	for rows.Next() {
		var course class_domain.Course
		_ = rows.Scan(&course.ID, &course.Name, &course.CreditHr, &course.Color)

		populateCourse(courseRepo.db, &course)
		courses = append(courses, course)
	}

	return courses, nil
}

func (courseRepo PsqlClassRepoAdapter) UpdateCourse(course class_domain.Course) (class_domain.Course, error) {
	var updatedCourse class_domain.Course
	err := courseRepo.db.QueryRow("UPDATE public.courses SET name = $1, credit_hr = $2, color = $3 WHERE course_id = $4 RETURNING *", course.Name, course.CreditHr, course.Color, course.ID).Scan(&updatedCourse.ID, &updatedCourse.Name, &updatedCourse.CreditHr, &updatedCourse.Color)
	if err != nil {
		return updatedCourse, err
	}

	populateCourse(courseRepo.db, &updatedCourse)
	return updatedCourse, nil
}

func (courseRepo PsqlClassRepoAdapter) DeleteCourse(id string) error {
	_, err := courseRepo.db.Query("DELETE FROM public.courses WHERE course_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (courseRepo PsqlClassRepoAdapter) AddUserToCourse(courseId, userId string) (class_domain.Course, error) {
	_, err := courseRepo.db.Query("INSERT INTO public.instructor_courses (course_id, user_id) VALUES ($1,$2)", courseId, userId)

	if err != nil {
		return class_domain.Course{}, err
	}

	return courseRepo.FindCourse(courseId)
}

func (courseRepo PsqlClassRepoAdapter) RemoveUserFromCourse(courseId, userId string) (class_domain.Course, error) {
	_, err := courseRepo.db.Query("DELETE FROM public.instructor_courses WHERE course_id = $1 AND user_id = $2", courseId, userId)
	if err != nil {
		return class_domain.Course{}, err
	}

	return courseRepo.FindCourse(courseId)
}
