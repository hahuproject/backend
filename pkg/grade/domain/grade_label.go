package grade_domain

import "time"

type GradeLabel struct {
	ID        string    `json:"gradeLabelId"`
	Label     string    `json:"label"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	CreatedAt time.Time `json:"createdAt"`
}
