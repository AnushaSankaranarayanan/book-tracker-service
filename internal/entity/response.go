package entity

import "net/http"

type GenericResponse struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BookResponse struct {
	GenericResponse
	Book  *Book  `json:"book,omitempty"`
	Count int    `json:"count,omitempty"`
	Books []Book `json:"books,omitempty"`
}

type GroupByGenreResponse struct {
	GenericResponse
	Genres []BooksByGenre `json:"genres"`
}

func NewGenericResponse(code int, msg string) GenericResponse {
	return GenericResponse{Code: code, Status: http.StatusText(code), Message: msg}
}

func NewBookResponse(code int, msg string, book *Book, books []Book) BookResponse {
	return BookResponse{
		GenericResponse: GenericResponse{
			Code:    code,
			Status:  http.StatusText(code),
			Message: msg,
		},
		Book:  book,
		Books: books,
		Count: len(books),
	}
}
func NewGroupByGenreResponse(code int, msg string, genres []BooksByGenre) GroupByGenreResponse {
	return GroupByGenreResponse{
		GenericResponse: GenericResponse{
			Code:    code,
			Status:  http.StatusText(code),
			Message: msg,
		},
		Genres: genres,
	}
}
