package library_psql_repo

import (
	"database/sql"
	"log"

	libraryDomain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/domain"
	libraryRepo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/repo"
	library_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/utils"
)

type LibraryPsqlRepoAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewLibraryPsqlRepoAdapter(log *log.Logger, db *sql.DB) (libraryRepo.LibraryRepoPort, error) {
	err := db.Ping()
	if err != nil {
		return &LibraryPsqlRepoAdapter{}, err
	}
	return &LibraryPsqlRepoAdapter{log: log, db: db}, nil
}

func populateBook(db *sql.DB, book *libraryDomain.Book) {
	_ = db.QueryRow(`SELECT
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.users
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.users.user_id = $1`, book.Uploader.ID).Scan(
		&book.Uploader.ID, &book.Uploader.FirstName, &book.Uploader.LastName, &book.Uploader.Email, &book.Uploader.Phone, &book.Uploader.Username, &book.Uploader.ProfilePic, &book.Uploader.Verified, &book.Uploader.Type,
		&book.Uploader.Address.ID, &book.Uploader.Address.Country, &book.Uploader.Address.Region, &book.Uploader.Address.City, &book.Uploader.Address.SubCity, &book.Uploader.Address.Woreda, &book.Uploader.Address.HouseNo)

	_ = db.QueryRow(`SELECT
			public.courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color
		FROM public.courses
		WHERE public.courses.course_id = $1`, book.Course.ID).Scan(
		&book.Course.ID, &book.Course.Name, &book.Course.CreditHr, &book.Course.Color)
}

func (libraryRepo LibraryPsqlRepoAdapter) StoreBook(book libraryDomain.Book) (libraryDomain.Book, error) {
	var addedBook libraryDomain.Book
	err := libraryRepo.db.QueryRow(`
	INSERT 
	INTO public.books 
	(title, author, description, cover, file, uploader, course_id) 
	VALUES ($1,$2,$3,$4,$5, $6, $7) RETURNING book_id, title, author, description, uploader, cover, file`,
		book.Title, book.Author, book.Description, book.Cover, book.File, book.Uploader.ID, sql.NullString{Valid: book.Course.ID != "", String: book.Course.ID}).
		Scan(&addedBook.ID, &addedBook.Title, &addedBook.Author, &addedBook.Description, &addedBook.Uploader.ID, &addedBook.Cover, &addedBook.File)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "title"` {
			libraryRepo.log.Println(err)
			return libraryDomain.Book{}, library_utils.ErrBookTitleTaken
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "cover"` {
			libraryRepo.log.Println(err)
			return libraryDomain.Book{}, library_utils.ErrBookCoverTaken
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "file"` {
			libraryRepo.log.Println(err)
			return libraryDomain.Book{}, library_utils.ErrBookFileTaken
		}
		libraryRepo.log.Println(err)
		return addedBook, library_utils.ErrFailedToStoreBook
	}
	return addedBook, nil
}
func (libraryRepo LibraryPsqlRepoAdapter) FindBooks() ([]libraryDomain.Book, error) {
	var books []libraryDomain.Book = make([]libraryDomain.Book, 0)

	rows, err := libraryRepo.db.Query(`
	SELECT 
	public.books.book_id, public.books.title, public.books.author, public.books.description, public.books.uploader, public.books.cover, public.books.file, public.books.rating, public.books.course_id
	FROM public.books`)
	if err != nil {
		return []libraryDomain.Book{}, library_utils.ErrNoBooksFound
	}

	for rows.Next() {
		var book libraryDomain.Book
		var _sqlNull sql.NullString
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Uploader.ID, &book.Cover, &book.File, &book.Rating, &_sqlNull)
		if err != nil {
			libraryRepo.log.Println(err)
		}
		if _sqlNull.Valid {
			book.Course.ID = _sqlNull.String
		}
		populateBook(libraryRepo.db, &book)
		books = append(books, book)
	}

	return books, nil
}
func (libraryRepo LibraryPsqlRepoAdapter) FindBookById(id string) (libraryDomain.Book, error) {
	var book libraryDomain.Book
	err := libraryRepo.db.QueryRow("SELECT * FROM public.books WHERE book_id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Uploader, &book.Cover, &book.File)
	if err != nil {
		return book, err
	}
	populateBook(libraryRepo.db, &book)
	return book, nil
}
func (libraryRepo LibraryPsqlRepoAdapter) UpdateBook(book libraryDomain.Book) (libraryDomain.Book, error) {
	return libraryDomain.Book{}, nil
}
func (libraryRepo LibraryPsqlRepoAdapter) UpdateBookRating(bookId string, rating float32) error {
	_, err := libraryRepo.db.Query(`UPDATE public.books SET rating = $1 WHERE book_id = $2`, rating, bookId)
	if err != nil {
		return err
	}
	return nil
}
func (libraryRepo LibraryPsqlRepoAdapter) DeleteBook(id string) error {
	return nil
}
