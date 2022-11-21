package announcement_service

import (
	"log"

	announcement_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/domain"
	announcement_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/repo"
	announcement_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/utils"
)

type AnnouncementServicePort interface {
	GetAnnouncements(token string) ([]announcement_domain.Announcement, error)
	GetAnnouncementsBySection(token, sectionId string) ([]announcement_domain.Announcement, error)
	AddAnnouncement(token string, announcement announcement_domain.Announcement) (announcement_domain.Announcement, error)
	UpdateAnnouncement(token string, announcement announcement_domain.Announcement) (announcement_domain.Announcement, error)
	DeleteAnnouncement(token, id string) error
}

type AnnouncementServiceAdapter struct {
	log  *log.Logger
	repo announcement_repo.AnnouncementRepoPort
}

func NewAnnouncementServiceAdapter(log *log.Logger, repo announcement_repo.AnnouncementRepoPort) AnnouncementServicePort {
	return &AnnouncementServiceAdapter{log: log, repo: repo}
}

func (announcementService AnnouncementServiceAdapter) GetAnnouncements(token string) ([]announcement_domain.Announcement, error) {

	var announcements []announcement_domain.Announcement = make([]announcement_domain.Announcement, 0)

	user, err := announcement_utils.CheckAuth(token)
	if err != nil {
		return announcements, err
	}

	if user.Type == "SUPER_ADMIN" || user.Type == "ADMIN" {
		return announcementService.repo.FindAnnouncements()
	}

	if user.Type == "INSTRUCTOR" || user.Type == "SUPERVISOR" || user.Type == "REGISTRY_OFFICER" || user.Type == "DEPARTMENT_HEAD" {
		return announcementService.repo.FindUserAnnonucements(user.ID)
	}

	if user.Type == "STUDENT" {
		return announcementService.repo.FindStudentAnnouncements(user.ID)
	}

	return announcements, announcement_utils.ErrUnauthorized
}

func (announcementService AnnouncementServiceAdapter) GetAnnouncementsBySection(token, sectionId string) ([]announcement_domain.Announcement, error) {
	var announcements []announcement_domain.Announcement = make([]announcement_domain.Announcement, 0)

	_, err := announcement_utils.CheckAuth(token)
	if err != nil {
		return announcements, err
	}

	return announcementService.repo.FindAnnouncementsBySection(sectionId)
}

func (announcementService AnnouncementServiceAdapter) AddAnnouncement(token string, announcement announcement_domain.Announcement) (announcement_domain.Announcement, error) {

	user, err := announcement_utils.CheckAuth(token)
	if err != nil {
		return announcement_domain.Announcement{}, err
	}

	announcement.PostedBy = user

	return announcementService.repo.StoreAnnouncement(announcement)
}
func (announcementService AnnouncementServiceAdapter) UpdateAnnouncement(token string, announcement announcement_domain.Announcement) (announcement_domain.Announcement, error) {
	user, err := announcement_utils.CheckAuth(token)
	if err != nil {
		return announcement_domain.Announcement{}, err
	}

	if user.ID != announcement.PostedBy.ID {
		return announcement_domain.Announcement{}, announcement_utils.ErrUnauthorized
	}

	return announcementService.repo.UpdateAnnouncement(announcement)
}
func (announcementService AnnouncementServiceAdapter) DeleteAnnouncement(token, id string) error {

	user, err := announcement_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	return announcementService.repo.DeleteAnnouncement(user.ID, id)
}
