package grade_psql_repo

import (
	"database/sql"
	"log"

	grade_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/domain"
	grade_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/repo"
)

type GradePsqlRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewGradePsqlRepoAdapter(log *log.Logger, db *sql.DB) (grade_repo.GradeRepoPort, error) {
	if err := db.Ping(); err != nil {
		return GradePsqlRepoAdapter{}, err
	}

	return GradePsqlRepoAdapter{log: log, db: db}, nil
}

func (repo GradePsqlRepoAdapter) StoreGrade(grade grade_domain.Grade) (grade_domain.Grade, error) {
	var addedGrade grade_domain.Grade

	// print("grade repo 0")

	var _ass sql.NullFloat64
	var _mid sql.NullFloat64
	var _final sql.NullFloat64
	var _lab sql.NullFloat64

	err := repo.db.QueryRow(`
	INSERT INTO public.grades
	(assessment, mid, "final", lab, course_id, user_id)
	VALUES ($1,$2,$3,$4, $5, $6)
	ON CONFLICT (course_id, user_id)
	DO UPDATE SET assessment = $1, mid = $2, "final" = $3, lab = $4
	WHERE public.grades.course_id = $5 AND public.grades.user_id = $6
	RETURNING grade_id, assessment,mid,final, lab, course_id, user_id, can_review, review_requested, submitted
	`, sql.NullFloat64{Valid: grade.Assessment > 0, Float64: grade.Assessment}, sql.NullFloat64{Valid: grade.Mid > 0, Float64: grade.Mid}, sql.NullFloat64{Valid: grade.Final > 0, Float64: grade.Final}, sql.NullFloat64{Valid: grade.Lab > 0, Float64: grade.Lab}, grade.Course.ID, grade.User.ID).Scan(
		&addedGrade.ID, &_ass, &_mid, &_final, &_lab, &addedGrade.Course.ID, &addedGrade.User.ID, &addedGrade.CanReview, &addedGrade.ReviewRequested, &addedGrade.Sunmitted,
	)

	if err != nil {
		return addedGrade, err
	}

	if _ass.Valid {
		addedGrade.Assessment = _ass.Float64
	}

	if _mid.Valid {
		addedGrade.Mid = _mid.Float64
	}

	if _final.Valid {
		addedGrade.Final = _final.Float64
	}

	if _lab.Valid {
		addedGrade.Lab = _lab.Float64
	}

	return addedGrade, nil
}

func (repo GradePsqlRepoAdapter) FindGrades() ([]grade_domain.Grade, error) {
	var grades []grade_domain.Grade = make([]grade_domain.Grade, 0)

	rows, err := repo.db.Query(`
	SELECT
	public.grades.grade_id, public.grades.assessment,public.grades.mid,public.grades.final,public.grades.lab, public.grades.course_id, public.grades.user_id, public.grades.can_review, public.grades.review_requested, public.grades.submitted,
	public.courses.name, public.courses.credit_hr,
	public.users.first_name, public.users.last_name
	FROM public.grades
	INNER JOIN public.courses ON public.courses.course_id = public.grades.course_id
	INNER JOIN public.users ON public.users.user_id = public.grades.user_id
	WHERE submitted = $1
	`, true)
	if err != nil {
		return grades, err
	}

	for rows.Next() {
		var grade grade_domain.Grade
		var _ass = sql.NullFloat64{}
		var _mid = sql.NullFloat64{}
		var _final = sql.NullFloat64{}
		var _lab = sql.NullFloat64{}
		rows.Scan(&grade.ID, &_ass, &_mid, &_final, &_lab, &grade.Course.ID, &grade.User.ID, &grade.CanReview, &grade.ReviewRequested, &grade.Sunmitted,
			&grade.Course.Name, &grade.Course.CreditHr,
			&grade.User.FirstName, &grade.User.LastName)

		if _ass.Valid {
			grade.Assessment = _ass.Float64
		}
		if _mid.Valid {
			grade.Mid = _mid.Float64
		}
		if _final.Valid {
			grade.Final = _final.Float64
		}
		if _lab.Valid {
			grade.Lab = _lab.Float64
		}

		grades = append(grades, grade)
	}

	return grades, nil
}

func (repo GradePsqlRepoAdapter) FindGradesByInstructor(instructorId string) ([]grade_domain.Grade, error) {
	var grades []grade_domain.Grade = make([]grade_domain.Grade, 0)

	rows, err := repo.db.Query(`
	SELECT
	public.grades.grade_id, public.grades.assessment,public.grades.mid,public.grades.final,public.grades.lab, public.grades.course_id, public.grades.user_id, public.grades.can_review, public.grades.review_requested, public.grades.submitted
	FROM public.grades
	INNER JOIN public.instructor_courses ON public.instructor_courses.course_id = public.grades.course_id
	WHERE public.instructor_courses.user_id = $1
	`, instructorId)
	if err != nil {
		return grades, err
	}

	for rows.Next() {
		var grade grade_domain.Grade
		var _ass = sql.NullFloat64{}
		var _mid = sql.NullFloat64{}
		var _final = sql.NullFloat64{}
		var _lab = sql.NullFloat64{}
		rows.Scan(&grade.ID, &_ass, &_mid, &_final, &_lab, &grade.Course.ID, &grade.User.ID, &grade.CanReview, &grade.ReviewRequested, &grade.Sunmitted)

		if _ass.Valid {
			grade.Assessment = _ass.Float64
		}
		if _mid.Valid {
			grade.Mid = _mid.Float64
		}
		if _final.Valid {
			grade.Final = _final.Float64
		}
		if _lab.Valid {
			grade.Lab = _lab.Float64
		}

		grades = append(grades, grade)
	}

	return grades, nil
}

func (repo GradePsqlRepoAdapter) FindGradesByUser(userId string) ([]grade_domain.Grade, error) {
	var grades []grade_domain.Grade = make([]grade_domain.Grade, 0)

	// repo.log.Println("grades[0].Course.ID")
	rows, err := repo.db.Query(`
	SELECT
	grade_id, assessment, mid, final, public.grades.user_id, can_review, review_requested, submitted,
	public.users.first_name, public.users.last_name, 
	public.sections.name, public.classes.name,
	public.courses.course_id, public.courses.name, public.courses.credit_hr
	FROM public.grades
	INNER JOIN public.users ON public.users.user_id = public.grades.user_id
	INNER JOIN public.courses ON public.courses.course_id = public.grades.course_id
	INNER JOIN public.students ON public.students.user_id = public.grades.user_id
	LEFT JOIN public.sections ON public.sections.section_id = public.students.section_id
	LEFT JOIN public.classes ON public.classes.class_id = public.sections.class_id
	WHERE public.grades.user_id = $1
	`, userId)
	if err != nil {
		repo.log.Println(err)
		return grades, err
	}

	for rows.Next() {
		var grade grade_domain.Grade
		var _ass = sql.NullFloat64{}
		var _mid = sql.NullFloat64{}
		var _final = sql.NullFloat64{}

		var _sectionName = sql.NullString{}
		var _className = sql.NullString{}

		err = rows.Scan(
			&grade.ID, &_ass, &_mid, &_final, &grade.User.ID, &grade.CanReview, &grade.ReviewRequested, &grade.Sunmitted,
			&grade.User.FirstName, &grade.User.LastName,
			&_sectionName, &_className,
			&grade.Course.ID, &grade.Course.Name, &grade.Course.CreditHr)
		if err != nil {
			repo.log.Println(err)
		}

		if _ass.Valid {
			grade.Assessment = _ass.Float64
		}
		if _mid.Valid {
			grade.Mid = _mid.Float64
		}
		if _final.Valid {
			grade.Final = _final.Float64
		}

		repo.log.Println(grade.Course.ID)
		grades = append(grades, grade)
	}

	repo.log.Println(grades[0].Course.ID)

	return grades, nil
}

func (repo GradePsqlRepoAdapter) SubmitGrade(userId, courseId string) error {
	_, err := repo.db.Query(`
	UPDATE public.grades
	SET submitted = $3
	WHERE user_id = $1 AND course_id = $2
	`, userId, courseId, true)

	return err
}

func (repo GradePsqlRepoAdapter) RequestGradeReview(userId, courseId string) error {

	_, err := repo.db.Query(`
	UPDATE public.grades
	SET review_requested = $3
	WHERE user_id = $1 AND course_id = $2
	`, userId, courseId, true)

	return err
}

func (repo GradePsqlRepoAdapter) ApproveGradeReview(userId, courseId string) error {

	_, err := repo.db.Query(`
	UPDATE
	public.grades
	SET can_review = $3, review_requested = $4, submitted = $5
	WHERE user_id = $1 AND course_id = $2
	`, userId, courseId, true, false, false)

	return err
}

func (repo GradePsqlRepoAdapter) RejectGradeReview(userId, courseId string) error {

	_, err := repo.db.Query(`
	UPDATE
	public.grades
	SET can_review = $3, review_requested = $4
	WHERE user_id = $1 AND course_id = $2
	`, userId, courseId, false, false)

	return err
}

//Grade Label
func (repo GradePsqlRepoAdapter) StoreGradeLabel(gradeLabel grade_domain.GradeLabel) (grade_domain.GradeLabel, error) {

	var addedGradeLabel grade_domain.GradeLabel

	err := repo.db.QueryRow(`
	INSERT INTO 
	public.grade_labels 
	("label", "min", "max") 
	VALUES ($1,$2,$3)
	RETURNING public.grade_labels.grade_label_id, public.grade_labels.label, public.grade_labels.min, public.grade_labels.max`,
		gradeLabel.Label, gradeLabel.Min, gradeLabel.Max).
		Scan(&addedGradeLabel.ID, &addedGradeLabel.Label, &addedGradeLabel.Min, &addedGradeLabel.Max)

	return addedGradeLabel, err
}

func (repo GradePsqlRepoAdapter) FindGradeLabels() ([]grade_domain.GradeLabel, error) {

	var gradeLabels []grade_domain.GradeLabel = make([]grade_domain.GradeLabel, 0)

	rows, err := repo.db.Query(`
	SELECT public.grade_labels.grade_label_id, public.grade_labels.label, public.grade_labels.min, public.grade_labels.max
	FROM public.grade_labels`)
	if err != nil {
		return gradeLabels, err
	}

	for rows.Next() {
		var gradeLabel grade_domain.GradeLabel
		rows.Scan(&gradeLabel.ID, &gradeLabel.Label, &gradeLabel.Min, &gradeLabel.Max)

		gradeLabels = append(gradeLabels, gradeLabel)
	}

	return gradeLabels, nil

}

func (repo GradePsqlRepoAdapter) RemoveGradeLabels(gradeLabelId string) error {
	_, err := repo.db.Exec(`DELETE FROM public.grade_labels WHERE grade_label_id = $1`, gradeLabelId)
	return err
}
