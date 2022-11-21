package annoucement_repo

import annoucement_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/domain"

type AnnouncementRepoPort interface {
	FindAnnouncements() ([]annoucement_domain.Announcement, error)
	FindUserAnnonucements(userId string) ([]annoucement_domain.Announcement, error)
	FindAnnouncementsBySection(sectionId string) ([]annoucement_domain.Announcement, error)
	FindStudentAnnouncements(studentId string) ([]annoucement_domain.Announcement, error)
	StoreAnnouncement(announcement annoucement_domain.Announcement) (annoucement_domain.Announcement, error)
	UpdateAnnouncement(announcement annoucement_domain.Announcement) (annoucement_domain.Announcement, error)
	DeleteAnnouncement(userId, id string) error
}
