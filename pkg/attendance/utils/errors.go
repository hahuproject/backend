package attendance_utils

import "errors"

var (
	//Book fields
	ErrInvalidBookTitle    = errors.New("invalid book title")
	ErrInvalidBookAuthor   = errors.New("invalid book author")
	ErrInvalidBookUploader = errors.New("invalid book uploader")
	ErrInvalidBookFile     = errors.New("invalid book file")

	//Auth
	ErrNotAuthorized = errors.New("not authorized for the operation")

	//DB
	ErrFailedToStoreBook = errors.New("failed to save book to database")
	ErrBookTitleTaken    = errors.New("book with the same title found")
	ErrBookCoverTaken    = errors.New("book with the same cover found")
	ErrBookFileTaken     = errors.New("book with the same file found")
	ErrNoBooksFound      = errors.New("no registered books found")

	ErrBookNotFound = errors.New("book could not be found")
)
