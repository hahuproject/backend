package library_handler

import "net/http"

type LibraryHandlerPort interface {
	GetAddBook(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetBook(w http.ResponseWriter, r *http.Request)
	GetUpdateBook(w http.ResponseWriter, r *http.Request)
	GetUpdateBookRating(w http.ResponseWriter, r *http.Request)
	GetDeleteBook(w http.ResponseWriter, r *http.Request)
}
