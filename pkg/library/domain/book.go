package library_domain

import (
	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	class_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/domain"
)

type Book struct {
	ID          string              `json:"bookId"`
	Title       string              `json:"title"`
	Author      string              `json:"author"`
	Description string              `json:"description"`
	Cover       string              `json:"cover"`
	File        string              `json:"file"`
	Uploader    auth_domain.User    `json:"uploader"`
	Rating      float32             `json:"rating"`
	Course      class_domain.Course `json:"course"`
}
