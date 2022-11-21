package class_domain

type Stream struct {
	ID         string     `json:"streamId"`
	Name       string     `json:"name"`
	Department Department `json:"department"`
}
