package library_repo

import (
	serviceDomain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/domain"
)

type LibraryRepoPort interface {
	StoreBook(book serviceDomain.Book) (serviceDomain.Book, error)
	FindBooks() ([]serviceDomain.Book, error)
	FindBookById(id string) (serviceDomain.Book, error)
	UpdateBook(book serviceDomain.Book) (serviceDomain.Book, error)
	UpdateBookRating(bookId string, rating float32) error
	DeleteBook(id string) error
}
