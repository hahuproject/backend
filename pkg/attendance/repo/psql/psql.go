package attendance_psql_repo

import (
	"database/sql"
	"log"

	attendance_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/domain"
	attendance_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/repo"
)

type AttendancePsqlRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewAttendancePsqlRepoAdapter(log *log.Logger, db *sql.DB) (attendance_repo.AttendanceRepoPort, error) {
	if err := db.Ping(); err != nil {
		return AttendancePsqlRepoAdapter{}, err
	}

	return AttendancePsqlRepoAdapter{log: log, db: db}, nil
}

func (repo AttendancePsqlRepoAdapter) StoreAttendance(attendances []attendance_domain.Attendance, userId string) ([]attendance_domain.Attendance, error) {

	var addedAttendances []attendance_domain.Attendance

	tx, err := repo.db.Begin()
	if err != nil {
		return addedAttendances, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, obj := range attendances {
		if _, err := tx.Exec(`
		INSERT INTO public.attendances 
		(user_id, status, session_id)
		VALUES ($1,$2,$3)
		ON CONFLICT (user_id, session_id)
		DO UPDATE SET status = $2
		WHERE public.attendances.user_id = $1 AND public.attendances.session_id = $3
		`, obj.User.ID, obj.Status, obj.Session.ID); err != nil {
			repo.log.Println("Failed to insert attendance")
			repo.log.Println(err)
			tx.Rollback()
			return addedAttendances, err
		}
	}

	if _, err := tx.Exec(`
		INSERT INTO public.attendances 
		(user_id, status, session_id)
		VALUES ($1,$2,$3)
		ON CONFLICT (user_id, session_id)
		DO UPDATE SET status = $2
		WHERE public.attendances.user_id = $1 AND public.attendances.session_id = $3
		`, userId, "PRESENT", attendances[0].Session.ID); err != nil {
		repo.log.Println("Failed to insert attendance")
		repo.log.Println(err)
		tx.Rollback()
		return addedAttendances, err
	}

	err = tx.Commit()
	if err != nil {
		return addedAttendances, err
	}

	return repo.FindAttendancesBySession(attendances[0].Session.ID)

}

// func (repo AttendancePsqlRepoAdapter) FindAttendances() ([]attendance_domain.Attendance, error) {

// }
func (repo AttendancePsqlRepoAdapter) FindAttendancesByUser(userId string) ([]attendance_domain.Attendance, error) {

	var attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)

	rows, err := repo.db.Query(`
	SELECT 
		attendance_id, public.attendances.user_id, status, public.attendances.session_id, created_at,
		course_id
	FROM public.attendances
	INNER JOIN public.sessions ON public.sessions.session_id = public.attendances.session_id
	INNER JOIN public.instructor_courses ON public.instructor_courses.course_user_id = public.sessions.course_user_id
	WHERE public.attendances.user_id = $1
	`, userId)

	if err != nil {
		return attendances, err
	}

	for rows.Next() {
		var attendance attendance_domain.Attendance
		rows.Scan(&attendance.ID, &attendance.User.ID, &attendance.Status, &attendance.Session.ID, &attendance.CreatedAt, &attendance.Session.Course.ID)

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (repo AttendancePsqlRepoAdapter) FindAttendancesByInstructor(instructorId string) ([]attendance_domain.Attendance, error) {
	var attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)

	rows, err := repo.db.Query(`
	SELECT 
	public.attendances.attendance_id, public.attendances.user_id, public.attendances.status, public.attendances.session_id, public.attendances.created_at
	FROM public.attendances
		INNER JOIN public.sessions on public.sessions.session_id = public.attendances.session_id
		INNER JOIN public.instructor_courses on public.instructor_courses.course_user_id = public.sessions.course_user_id
	WHERE public.instructor_courses.user_id = $1 AND public.attendances.user_id != $1
	`, instructorId)

	if err != nil {
		return attendances, err
	}

	for rows.Next() {
		var attendance attendance_domain.Attendance
		rows.Scan(&attendance.ID, &attendance.User.ID, &attendance.Status, &attendance.Session.ID, &attendance.CreatedAt)

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (repo AttendancePsqlRepoAdapter) FindAttendancesBySession(sessionId string) ([]attendance_domain.Attendance, error) {
	var attendances []attendance_domain.Attendance = make([]attendance_domain.Attendance, 0)

	rows, err := repo.db.Query(`
	SELECT 
		attendance_id, user_id, status, session_id, created_at
	FROM public.attendances
	WHERE public.attendances.session_id = $1
	`, sessionId)

	if err != nil {
		return attendances, err
	}

	for rows.Next() {
		var attendance attendance_domain.Attendance
		rows.Scan(&attendance.ID, &attendance.User.ID, &attendance.Status, &attendance.Session.ID, &attendance.CreatedAt)

		attendances = append(attendances, attendance)
	}

	return attendances, nil
}
