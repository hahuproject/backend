package class_psql_repo

import (
	"database/sql"
	"log"

	class_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/repo"
	classError "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/utils"
)

type PsqlClassRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewPsqlCourseRepoAdapter(log *log.Logger, db *sql.DB) (class_repo.ClassRepoPort, error) {
	err := db.Ping()
	if err != nil {
		return &PsqlClassRepoAdapter{}, classError.ErrDbPingFailed
	}
	return &PsqlClassRepoAdapter{log: log, db: db}, nil
}
