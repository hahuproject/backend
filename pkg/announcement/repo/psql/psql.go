package announcement_psql_repo

import (
	"database/sql"
	"log"

	annoucement_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/domain"
	annoucement_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/repo"
	auth_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/repo/psql"
	class_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/repo/psql"
)

type PsqlAnnouncementAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewPsqlAnnouncementAdapter(log *log.Logger, db *sql.DB) (annoucement_repo.AnnouncementRepoPort, error) {
	err := db.Ping()
	if err != nil {
		return &PsqlAnnouncementAdapter{}, err
	}

	return &PsqlAnnouncementAdapter{
		log: log,
		db:  db,
	}, nil
}

func populateAnnouncement(repo PsqlAnnouncementAdapter, announcement *annoucement_domain.Announcement) {
	//Posted By
	authRepo, err := auth_psql_repo.NewPsqlAuthRepositoryAdapter(repo.log, repo.db)
	if err == nil {
		announcement.PostedBy, _ = authRepo.FindUserById(announcement.PostedBy.ID)
	}
	//Section
	if announcement.Section.ID != "" {
		sectionRepo, err := class_psql_repo.NewPsqlCourseRepoAdapter(repo.log, repo.db)
		if err == nil {
			announcement.Section, _ = sectionRepo.FindSection(announcement.Section.ID)
		}
	}
}

func (annoucementRepo PsqlAnnouncementAdapter) FindAnnouncements() ([]annoucement_domain.Announcement, error) {
	var announcements []annoucement_domain.Announcement = make([]annoucement_domain.Announcement, 0)

	rows, err := annoucementRepo.db.Query(`SELECT announcement_id, title, message, created_at, section_id, posted_by FROM public.announcements`)
	if err != nil {
		return announcements, err
	}

	for rows.Next() {
		var announcement annoucement_domain.Announcement
		var nullSection sql.NullString
		rows.Scan(&announcement.ID, &announcement.Title, &announcement.Message, &announcement.CreatedAt, &nullSection, &announcement.PostedBy.ID)

		if nullSection.Valid {
			announcement.Section.ID = nullSection.String
		}

		populateAnnouncement(annoucementRepo, &announcement)

		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) FindUserAnnonucements(userId string) ([]annoucement_domain.Announcement, error) {
	var announcements []annoucement_domain.Announcement = make([]annoucement_domain.Announcement, 0)

	rows, err := annoucementRepo.db.Query(`SELECT announcement_id, title, message, created_at, section_id, posted_by FROM public.announcements WHERE section_id IS NULL OR posted_by = $1`, userId)
	if err != nil {
		return announcements, err
	}

	for rows.Next() {
		var announcement annoucement_domain.Announcement
		var nullSection sql.NullString
		rows.Scan(&announcement.ID, &announcement.Title, &announcement.Message, &announcement.CreatedAt, &nullSection, &announcement.PostedBy.ID)

		if nullSection.Valid {
			announcement.Section.ID = nullSection.String
		}

		populateAnnouncement(annoucementRepo, &announcement)

		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) FindStudentAnnouncements(studentId string) ([]annoucement_domain.Announcement, error) {
	var announcements []annoucement_domain.Announcement = make([]annoucement_domain.Announcement, 0)

	rows, err := annoucementRepo.db.Query(`
	SELECT public.announcements.announcement_id, public.announcements.title, public.announcements.message, public.announcements.created_at, public.announcements.section_id, public.announcements.posted_by 
FROM public.announcements 
LEFT JOIN public.student_sections ON public.student_sections.section_id = public.announcements.section_id
WHERE public.announcements.section_id IS NULL OR public.student_sections.student_id  = $1`, studentId)
	if err != nil {
		return announcements, err
	}

	for rows.Next() {
		var announcement annoucement_domain.Announcement
		var nullSection sql.NullString
		rows.Scan(&announcement.ID, &announcement.Title, &announcement.Message, &announcement.CreatedAt, &nullSection, &announcement.PostedBy.ID)

		if nullSection.Valid {
			announcement.Section.ID = nullSection.String
		}

		populateAnnouncement(annoucementRepo, &announcement)

		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) FindAnnouncementsBySection(sectionId string) ([]annoucement_domain.Announcement, error) {
	var announcements []annoucement_domain.Announcement = make([]annoucement_domain.Announcement, 0)

	rows, err := annoucementRepo.db.Query(`SELECT announcement_id, title, message, created_at, section_id, posted_by FROM public.announcements WHERE section_id = $1`, sectionId)
	if err != nil {
		return announcements, err
	}

	for rows.Next() {
		var announcement annoucement_domain.Announcement
		var nullSection sql.NullString
		rows.Scan(&announcement.ID, &announcement.Title, &announcement.Message, &announcement.CreatedAt, &nullSection, &announcement.PostedBy.ID)

		if nullSection.Valid {
			announcement.Section.ID = nullSection.String
		}

		populateAnnouncement(annoucementRepo, &announcement)

		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) StoreAnnouncement(announcement annoucement_domain.Announcement) (annoucement_domain.Announcement, error) {

	err := annoucementRepo.db.QueryRow(`INSERT INTO public.announcements (title, message, posted_by, section_id) VALUES($1,$2,$3,$4) RETURNING announcement_id, created_at`, announcement.Title, announcement.Message, announcement.PostedBy.ID, sql.NullString{String: announcement.Section.ID, Valid: announcement.Section.ID != ""}).Scan(&announcement.ID, &announcement.CreatedAt)
	if err != nil {
		return announcement, err
	}

	return announcement, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) UpdateAnnouncement(announcement annoucement_domain.Announcement) (annoucement_domain.Announcement, error) {

	_, err := annoucementRepo.db.Query(`UPDATE public.announcements SET title = $1, message = $2 WHERE announcement_id = $3`, announcement.Title, announcement.Message, announcement.ID)
	if err != nil {
		return announcement, err
	}

	return announcement, nil
}

func (annoucementRepo PsqlAnnouncementAdapter) DeleteAnnouncement(userId, announcementId string) error {

	_, err := annoucementRepo.db.Query(`DELETE FROM public.announcements WHERE announcement_id = $1 AND posted_by = $2`, announcementId, userId)
	if err != nil {
		return err
	}

	return nil
}
