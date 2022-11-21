package assignment_psql_repo

import (
	"database/sql"
	"log"

	assignment_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/domain"
	assignment_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/repo"
)

type AssignmentPsqlRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewAssignmentPsqlRepoAdapter(log *log.Logger, db *sql.DB) (assignment_repo.AssignmentRepoPort, error) {
	err := db.Ping()
	if err != nil {
		return AssignmentPsqlRepoAdapter{}, err
	}

	return AssignmentPsqlRepoAdapter{log: log, db: db}, nil
}

func (repo AssignmentPsqlRepoAdapter) FindAssignments() ([]assignment_domain.Assignment, error) {

	var _assignments []assignment_domain.Assignment = make([]assignment_domain.Assignment, 0)

	rows, err := repo.db.Query(`
	SELECT 
		public.assignments.assignment_id, public.assignments.title, public.assignments.remark, public.assignments.attachment,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.sections.section_id, public.sections.name, public.sections.class_id
	FROM public.assignments
	INNER JOIN public.courses ON public.courses.course_id = public.assignments.course_id
	INNER JOIN public.sections ON public.sections.section_id = public.assignments.section_id
	`)

	if err != nil {
		return _assignments, err
	}

	for rows.Next() {
		var _assignment assignment_domain.Assignment
		rows.Scan(
			&_assignment.ID, &_assignment.Title, &_assignment.Remark, &_assignment.Attachment,
			&_assignment.Course.ID, &_assignment.Course.Name, &_assignment.Course.CreditHr, &_assignment.Course.Color,
			&_assignment.Section.ID, &_assignment.Section.Name, &_assignment.Section.Class.ID,
		)

		_assignments = append(_assignments, _assignment)
	}

	return _assignments, nil

}

func (repo AssignmentPsqlRepoAdapter) FindAssignmentById(id string) (assignment_domain.Assignment, error) {
	var _assignment assignment_domain.Assignment

	err := repo.db.QueryRow(`
	SELECT 
		public.assignments.assignment_id, public.assignments.title, public.assignments.remark, public.assignments.attachment,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.sections.section_id, public.sections.name, public.sections.class_id
	FROM public.assignments
	INNER JOIN public.courses ON public.courses.course_id = public.assignments.course_id
	INNER JOIN public.sections ON public.sections.section_id = public.assignments.section_id WHERE public.assignments.assignment_id = $1`, id).Scan(
		&_assignment.ID, &_assignment.Title, &_assignment.Remark, &_assignment.Attachment,
		&_assignment.Course.ID, &_assignment.Course.Name, &_assignment.Course.CreditHr, &_assignment.Course.Color,
		&_assignment.Section.ID, &_assignment.Section.Name, &_assignment.Section.Class.ID,
	)

	if err != nil {
		return _assignment, err
	}

	return _assignment, nil
}

func (repo AssignmentPsqlRepoAdapter) StoreAssignment(assignment assignment_domain.Assignment) (assignment_domain.Assignment, error) {

	var _addedAssignment assignment_domain.Assignment

	err := repo.db.QueryRow(`
	INSERT INTO public.assignments
	(title, remark, attachment, section_id, course_id)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING assignment_id`,
		assignment.Title, assignment.Remark, assignment.Attachment, assignment.Section.ID, assignment.Course.ID).Scan(&_addedAssignment.ID)
	if err != nil {
		return _addedAssignment, err
	}

	_addedAssignment, err = repo.FindAssignmentById(_addedAssignment.ID)
	if err != nil {
		return _addedAssignment, err
	}

	return _addedAssignment, nil
}

func (repo AssignmentPsqlRepoAdapter) UpdateAssignment(assignment assignment_domain.Assignment) (assignment_domain.Assignment, error) {
	var _updatedAssignment assignment_domain.Assignment

	err := repo.db.QueryRow(`
	UPDATE public.assignments
	SET title = $1, remark = $2, attachment = $3
	WHERE public.assignments.assignment_id = $4
	RETURNING public.assignments.assignment_id
	`, assignment.Title, assignment.Remark, assignment.Attachment, assignment.ID).Scan(&_updatedAssignment.ID)

	if err != nil {
		return _updatedAssignment, err
	}

	_updatedAssignment, err = repo.FindAssignmentById(_updatedAssignment.ID)
	if err != nil {
		return _updatedAssignment, err
	}

	return _updatedAssignment, nil
}

func (repo AssignmentPsqlRepoAdapter) DeleteAssignment(id string) error {
	_, err := repo.db.Query(`DELETE FROM public.assignments WHRER assignment_id = $1`, id)
	return err
}

func (repo AssignmentPsqlRepoAdapter) FindAssignmentsBySection(sectionId string) ([]assignment_domain.Assignment, error) {
	var _assignments []assignment_domain.Assignment = make([]assignment_domain.Assignment, 0)

	rows, err := repo.db.Query(`
	SELECT 
		public.assignments.assignment_id, public.assignments.title, public.assignments.remark, public.assignments.attachment,public.assignments.created_at,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.sections.section_id, public.sections.name, public.sections.class_id
	FROM public.assignments
	INNER JOIN public.courses ON public.courses.course_id = public.assignments.course_id
	INNER JOIN public.sections ON public.sections.section_id = public.assignments.section_id
	WHERE public.assignments.section_id = $1
	`, sectionId)

	if err != nil {
		return _assignments, err
	}

	for rows.Next() {
		var _assignment assignment_domain.Assignment
		rows.Scan(
			&_assignment.ID, &_assignment.Title, &_assignment.Remark, &_assignment.Attachment, &_assignment.CreatedAt,
			&_assignment.Course.ID, &_assignment.Course.Name, &_assignment.Course.CreditHr, &_assignment.Course.Color,
			&_assignment.Section.ID, &_assignment.Section.Name, &_assignment.Section.Class.ID,
		)

		_assignments = append(_assignments, _assignment)
	}

	return _assignments, nil
}

func (repo AssignmentPsqlRepoAdapter) FindAssignmentsByStudent(studentId string) ([]assignment_domain.Assignment, error) {
	var _assignments []assignment_domain.Assignment = make([]assignment_domain.Assignment, 0)

	rows, err := repo.db.Query(`
	SELECT 
		public.assignments.assignment_id, public.assignments.title, public.assignments.remark, public.assignments.attachment, public.assignments.created_at,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.sections.section_id, public.sections.name, public.sections.class_id
	FROM public.assignments
	INNER JOIN public.courses ON public.courses.course_id = public.assignments.course_id
	INNER JOIN public.sections ON public.sections.section_id = public.assignments.section_id
	INNER JOIN public.student_sections ON public.student_sections.section_id = public.assignments.section_id
		WHERE public.student_sections.student_id  = $1
	`, studentId)

	if err != nil {
		return _assignments, err
	}

	for rows.Next() {
		var _assignment assignment_domain.Assignment
		rows.Scan(
			&_assignment.ID, &_assignment.Title, &_assignment.Remark, &_assignment.Attachment, &_assignment.CreatedAt,
			&_assignment.Course.ID, &_assignment.Course.Name, &_assignment.Course.CreditHr, &_assignment.Course.Color,
			&_assignment.Section.ID, &_assignment.Section.Name, &_assignment.Section.Class.ID,
		)

		_assignments = append(_assignments, _assignment)
	}

	return _assignments, nil
}

func (repo AssignmentPsqlRepoAdapter) FindAssignmentsByInstructor(instructorId string) ([]assignment_domain.Assignment, error) {
	var _assignments []assignment_domain.Assignment = make([]assignment_domain.Assignment, 0)

	rows, err := repo.db.Query(`
	SELECT 
		public.assignments.assignment_id, public.assignments.title, public.assignments.remark, public.assignments.attachment, public.assignments.created_at,
		public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color,
		public.sections.section_id, public.sections.name, public.sections.class_id
	FROM public.assignments
	INNER JOIN public.courses ON public.courses.course_id = public.assignments.course_id
	INNER JOIN public.sections ON public.sections.section_id = public.assignments.section_id
	INNER JOIN public.instructor_courses ON public.instructor_courses.course_id = public.assignments.course_id
	WHERE public.instructor_courses.user_id = $1
	`, instructorId)

	if err != nil {
		return _assignments, err
	}

	for rows.Next() {
		var _assignment assignment_domain.Assignment
		rows.Scan(
			&_assignment.ID, &_assignment.Title, &_assignment.Remark, &_assignment.Attachment, &_assignment.CreatedAt,
			&_assignment.Course.ID, &_assignment.Course.Name, &_assignment.Course.CreditHr, &_assignment.Course.Color,
			&_assignment.Section.ID, &_assignment.Section.Name, &_assignment.Section.Class.ID,
		)

		_assignments = append(_assignments, _assignment)
	}

	return _assignments, nil
}
