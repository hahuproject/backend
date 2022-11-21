package library_rest_handler

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	library_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/domain"
	library_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/handler"
	library_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/service"
	library_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/utils"
)

const (
	BUCKET_NAME = "hahusmsbucket"
	REGION      = "us-east-1"
	PART_SIZE   = 5000000
	RETRIES     = 2
)

var (
	s3session *s3.S3
)

func initS3() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials("AKIATLRHCAJ7OTROBI5Q", "lGHhY0IXUHtoqxYy1BjQWEHH+EXxOGNqKBDq728D", ""),
	})))
}

func StoreFileFromFormFileToAWSS3(file multipart.File, size int64, title string) string {
	initS3()
	var storedFileAdd string = ""

	if file != nil {

		buffer := make([]byte, size)
		file.Read(buffer)

		expiryDate := time.Now().AddDate(0, 0, 1)

		createdRes, err := s3session.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
			Bucket:  aws.String(BUCKET_NAME),
			Key:     aws.String(title),
			Expires: &expiryDate,
		})

		if err != nil {
			log.Println("err in aws")
			log.Println(err)
			return storedFileAdd
		}

		var start, currentSize int
		var remaining int = int(size)
		var partNum = 1
		var completedParts []*s3.CompletedPart

		for start = 0; remaining != 0; start += PART_SIZE {

			if remaining < PART_SIZE {
				currentSize = remaining
			} else {
				currentSize = PART_SIZE
			}

			completed, err := Upload(createdRes, buffer[start:start+currentSize], partNum)
			if err != nil {
				_, err := s3session.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
					Bucket:   createdRes.Bucket,
					Key:      createdRes.Key,
					UploadId: createdRes.UploadId,
				})
				if err != nil {
					log.Println("err in aws 1")
					log.Println(err)
					return storedFileAdd
				}
				log.Println("err aws 2")
				log.Println(err)
				return storedFileAdd
			}

			remaining -= currentSize
			log.Println("Uploading")

			partNum++

			completedParts = append(completedParts, completed)

		}

		res, err := s3session.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
			Bucket:   createdRes.Bucket,
			Key:      createdRes.Key,
			UploadId: createdRes.UploadId,
			MultipartUpload: &s3.CompletedMultipartUpload{
				Parts: completedParts,
			},
		})

		if err != nil {
			log.Println("err final upload aws")
			log.Println(err)
			return storedFileAdd
		}

		storedFileAdd = *res.Location
	}

	defer file.Close()

	return storedFileAdd
}

func Upload(res *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int) (completedPart *s3.CompletedPart, err error) {
	var try int
	for try <= RETRIES {
		uploadRes, err := s3session.UploadPart(&s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        res.Bucket,
			Key:           res.Key,
			PartNumber:    aws.Int64(int64(partNum)),
			UploadId:      res.UploadId,
			ContentLength: aws.Int64(int64(len(fileBytes))),
		})

		if err != nil {
			log.Println("err aws upload")
			log.Println(err)

			if try == RETRIES {
				return nil, err
			} else {
				try++
			}
		} else {
			return &s3.CompletedPart{
				ETag:       uploadRes.ETag,
				PartNumber: aws.Int64(int64(partNum)),
			}, nil
		}

	}

	return nil, nil
}

type LibraryRestHandlerAdapter struct {
	log            *log.Logger
	libraryService library_service.LibraryServicePort
}

func NewLibraryRestHandlerAdapter(log *log.Logger, libraryService library_service.LibraryServicePort) library_handler.LibraryHandlerPort {
	return &LibraryRestHandlerAdapter{log: log, libraryService: libraryService}
}

func (libraryHandler LibraryRestHandlerAdapter) GetAddBook(w http.ResponseWriter, r *http.Request) {
	var book library_domain.Book

	r.ParseMultipartForm(32 << 20)

	//Get Book File
	bookFile, bookFileHandler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a valid file"))
		return
	}

	defer bookFile.Close()

	storedBookFile := StoreFileFromFormFileToAWSS3(bookFile, bookFileHandler.Size, r.FormValue("title")+"-"+r.FormValue("author")+"."+strings.Split(bookFileHandler.Filename, ".")[len(strings.Split(bookFileHandler.Filename, "."))-1])
	if storedBookFile == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to save book file"))
		return
	}

	//Get Book Cover
	bookCover, bookCoverHandler, err := r.FormFile("cover")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a valid cover"))
		return
	}

	defer bookCover.Close()

	storedBookCover := StoreFileFromFormFileToAWSS3(bookCover, bookCoverHandler.Size, r.FormValue("title")+"-"+r.FormValue("title")+"-"+r.FormValue("author")+"_cover"+"."+strings.Split(bookCoverHandler.Filename, ".")[len(strings.Split(bookCoverHandler.Filename, "."))-1])
	if storedBookCover == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to save book file"))
		return
	}

	book = library_domain.Book{
		Title:       r.FormValue("title"),
		Author:      r.FormValue("author"),
		Description: r.FormValue("description"),
		File:        storedBookFile,
		Cover:       storedBookCover,
	}

	book.Course.ID = r.FormValue("course")

	// libraryHandler.log.Println(book)

	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}
	addedBook, err := libraryHandler.libraryService.AddBook(token[1], book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledAddedBook, _ := json.Marshal(addedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedBook)
	// w.Write([]byte("OK"))

}
func (libraryHandler LibraryRestHandlerAdapter) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := libraryHandler.libraryService.GetBooks()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledBooks, _ := json.Marshal(books)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(marshaledBooks)
}
func (libraryHandler LibraryRestHandlerAdapter) GetBook(w http.ResponseWriter, r *http.Request) {
	book, err := libraryHandler.libraryService.GetBook(r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	marshaledBook, _ := json.Marshal(book)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(marshaledBook)

}
func (libraryHandler LibraryRestHandlerAdapter) GetUpdateBook(w http.ResponseWriter, r *http.Request) {
}
func (libraryHandler LibraryRestHandlerAdapter) GetUpdateBookRating(w http.ResponseWriter, r *http.Request) {
	token, err := library_utils.CheckBearerTokenFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var book library_domain.Book
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = libraryHandler.libraryService.UpdateBookRating(token, book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully rated a book"))

}
func (libraryHandler LibraryRestHandlerAdapter) GetDeleteBook(w http.ResponseWriter, r *http.Request) {
}
