//go:build fake

package service

import (
	"errors"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/consts"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/framework/database"

	"testing"
)

const (
	createBook        = "AddBook"
	updateBook        = "UpdateBook"
	getAllBooks       = "ListBooks"
	getBook           = "GetBook"
	groupBooksByGenre = "GroupBooksByGenre"
)

func TestService(t *testing.T) {
	testBook := entity.Book{ISBN: "TEST-ISBN", Title: "Test Book", Author: "Test Author"}

	tests := []struct {
		testName      string
		errorExpected error
		arg           string
		serviceMethod string
		errorFlag     string
		sortKey       string
	}{
		{
			"CreateBook: should pass",
			nil,
			"",
			createBook,
			"",
			"",
		},
		{
			"CreateBook: should fail (force db error)",
			errors.New("Upsert error:forced collection upsert error"),
			"",
			createBook,
			"error",
			"",
		},
		{
			"UpdateBook: should pass",
			nil,
			"",
			updateBook,
			"",
			"",
		},
		{
			"UpdateBook: should fail (Forced collection error)",
			errors.New("get error:forced collection error"),
			"",
			updateBook,
			"true",
			"",
		},
		{
			"UpdateBook: should fail (update-error)",
			errors.New("Upsert error:forced collection upsert error"),
			testBook.ISBN,
			updateBook,
			"update-error",
			"",
		},
		{
			"GetBooks: should pass",
			nil,
			"",
			getAllBooks,
			"",
			"",
		},
		{
			"GetBooks: should pass(sorted by title)",
			nil,
			"",
			getAllBooks,
			"",
			consts.Title,
		},
		{
			"GetBooks: should pass(sorted by status)",
			nil,
			"",
			getAllBooks,
			"",
			consts.Status,
		},
		{
			"GetBook: should fail(book not found)",
			errors.New("book with id TEST-ISBN not found"),
			testBook.ISBN,
			getBook,
			"not-found-error",
			"",
		},
		{
			"GetBook: should pass",
			nil,
			testBook.ISBN,
			getBook,
			"",
			"",
		},
		{
			"GroupBooksByGenre: should fail(force read error)",
			errors.New("GetAll query error:forced query error"),
			"",
			groupBooksByGenre,
			"query-error",
			"",
		},
		{
			"GroupBooksByGenre: should pass",
			nil,
			"",
			groupBooksByGenre,
			"",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			couchbaseStorage, _ := database.NewFakeCouchbaseStorage(test.errorFlag)
			bookService := NewBookTracker(couchbaseStorage)

			var err error
			switch test.serviceMethod {
			case createBook:
				err = bookService.AddBook(testBook)
			case updateBook:
				err = bookService.UpdateBook(testBook)
			case getAllBooks:
				_, err = bookService.ListBooks(test.sortKey)
			case getBook:
				_, err = bookService.GetBook(test.arg)
			case groupBooksByGenre:
				_, err = bookService.GroupBooksByGenre()
			}

			if err == nil && err != test.errorExpected {
				t.Errorf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", test.serviceMethod, err, test.errorExpected)
			}

			if test.errorExpected != nil && test.errorExpected.Error() != err.Error() {
				t.Errorf("Function (%s) assert (error type is different from expected) -  got (%s) wanted (%s)", test.serviceMethod, err.Error(), test.errorExpected.Error())
			}
		})
	}
}
