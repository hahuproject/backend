package library_service

import (
	"log"

	libraryDomain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/domain"
	library_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/repo"
	library_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/utils"
)

type LibraryServicePort interface {
	AddBook(token string, book libraryDomain.Book) (libraryDomain.Book, error)
	GetBooks() ([]libraryDomain.Book, error)
	GetBook(id string) (libraryDomain.Book, error)
	UpdateBook(book libraryDomain.Book) (libraryDomain.Book, error)
	UpdateBookRating(token string, book libraryDomain.Book) error
	DeleteBook(id string) error
}

type LibraryServiceAdapter struct {
	log  *log.Logger
	repo library_repo.LibraryRepoPort
}

func NewLibraryServiceAdapter(log *log.Logger, repo library_repo.LibraryRepoPort) LibraryServicePort {
	return &LibraryServiceAdapter{log: log, repo: repo}
}

func (libraryService LibraryServiceAdapter) AddBook(token string, book libraryDomain.Book) (libraryDomain.Book, error) {
	//Check data validty
	//check token
	if token == "" {
		libraryService.log.Println("No token found")
		return libraryDomain.Book{}, library_utils.ErrNotAuthorized
	}
	//check book title
	if book.Title == "" {
		return libraryDomain.Book{}, library_utils.ErrInvalidBookTitle
	}
	//check book author
	if book.Author == "" {
		return libraryDomain.Book{}, library_utils.ErrInvalidBookAuthor
	}
	//check book file
	if book.File == "" {
		return libraryDomain.Book{}, library_utils.ErrInvalidBookFile
	}
	//check book file
	if book.File == "" {
		return libraryDomain.Book{}, library_utils.ErrInvalidBookFile
	}

	user, err := library_utils.CheckAuth(token)
	if err != nil {
		return libraryDomain.Book{}, err
	}

	book.Uploader = user

	libraryService.log.Println(book)

	addedBook, err := libraryService.repo.StoreBook(book)
	if err != nil {
		return addedBook, err
	}

	return addedBook, nil
}
func (libraryService LibraryServiceAdapter) GetBooks() ([]libraryDomain.Book, error) {
	return libraryService.repo.FindBooks()
}
func (libraryService LibraryServiceAdapter) GetBook(id string) (libraryDomain.Book, error) {
	return libraryService.repo.FindBookById(id)
}
func (libraryService LibraryServiceAdapter) UpdateBook(book libraryDomain.Book) (libraryDomain.Book, error) {
	return libraryDomain.Book{}, nil
}
func (libraryService LibraryServiceAdapter) UpdateBookRating(token string, book libraryDomain.Book) error {

	_, err := library_utils.CheckAuth(token)
	if err != nil {
		return err
	}

	_oldBook, err := libraryService.repo.FindBookById(book.ID)
	if err != nil {
		return err
	}

	_newRating := _oldBook.Rating + (book.Rating / 1000)

	return libraryService.repo.UpdateBookRating(book.ID, _newRating)
}
func (libraryService LibraryServiceAdapter) DeleteBook(id string) error {
	return nil
}
