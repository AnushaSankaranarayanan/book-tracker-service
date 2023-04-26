//go:build fake
// +build fake

package webserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/framework/database"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	testFolderPath           = "../../../tests/"
	createBookHandler        = "CreateBook"
	getBooksHandler          = "GetBooks"
	getBookHandler           = "GetBook"
	updateBookHandler        = "UpdateBook"
	bookJsonFile             = "book.json"
	bookMissingFieldJsonFile = "book-missing-mandatory-field.json"
)

var (
	bookURL = "/api/v1/book"
)

func TestHandlers(t *testing.T) {
	crulTests := []struct {
		testName           string
		httpMethod         string
		errorFlag          string
		errorExpected      string
		statusCodeExpected int
		requestPayload     string
		handler            string
	}{
		{
			"AddBook Book: missing mandatory fields",
			http.MethodPost,
			"",
			"Key: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
			http.StatusBadRequest,
			bookMissingFieldJsonFile,
			createBookHandler,
		},
		{
			"AddBook Book: force DB error",
			http.MethodPost,
			"error",
			"failed to save book.Refer to logs for more details",
			http.StatusInternalServerError,
			bookJsonFile,
			createBookHandler,
		},
		{
			"AddBook Book: should pass",
			http.MethodPost,
			"",
			"",
			http.StatusOK,
			bookJsonFile,
			createBookHandler,
		},
		{
			"UpdateBook Book: missing mandatory fields",
			http.MethodPut,
			"",
			"Key: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
			http.StatusBadRequest,
			bookMissingFieldJsonFile,
			updateBookHandler,
		},
		{
			"UpdateBook Book: document not found error",
			http.MethodPut,
			"not-found-error",
			"book with id TEST-ISBN-1 not found",
			http.StatusNotFound,
			bookJsonFile,
			updateBookHandler,
		},
		{
			"UpdateBook Book: force DB error",
			http.MethodPut,
			"update-error",
			"operation failed.Refer to logs for more details",
			http.StatusInternalServerError,
			bookJsonFile,
			updateBookHandler,
		},
		{
			"UpdateBook Book: should pass",
			http.MethodPut,
			"",
			"",
			http.StatusOK,
			bookJsonFile,
			updateBookHandler,
		},
		{
			"Get Book: document not found error",
			http.MethodGet,
			"not-found-error",
			"book with id  not found",
			http.StatusNotFound,
			"",
			getBookHandler,
		},
		{
			"Get Book: force DB error",
			http.MethodGet,
			"error",
			"operation failed.Refer to logs for more details",
			http.StatusInternalServerError,
			"",
			getBookHandler,
		},
		{
			"Get Book: should pass",
			http.MethodGet,
			"",
			"",
			http.StatusOK,
			"",
			getBookHandler,
		},
		{
			"Get Books: force DB error",
			http.MethodGet,
			"query-error",
			"failed to get books.Refer to logs for more details",
			http.StatusInternalServerError,
			"",
			getBooksHandler,
		},
		{
			"Get Books: should pass",
			http.MethodGet,
			"",
			"",
			http.StatusOK,
			"",
			getBooksHandler,
		},
	}

	for _, test := range crulTests {
		t.Run(test.testName, func(t *testing.T) {
			//setup
			file, _ := os.ReadFile(testFolderPath + test.requestPayload)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(test.httpMethod, bookURL, bytes.NewBuffer(file))
			c, _ := gin.CreateTestContext(rr)
			c.Request = req

			cbStorage, _ := database.NewFakeCouchbaseStorage(test.errorFlag)
			bookSvc := service.NewBookTracker(cbStorage)
			server := NewServer(Services{BookTracker: bookSvc})

			// actual tests
			switch test.handler {
			case createBookHandler:
				server.AddBook(c)
			case updateBookHandler:
				server.UpdateBook(c)
			case getBookHandler:
				server.GetBook(c)
			case getBooksHandler:
				server.ListBooks(c)
			}

			//assertions
			if rr.Code != test.statusCodeExpected {
				t.Errorf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", test.handler, rr.Code, test.statusCodeExpected)
			}

			//check error messages for failures
			if rr.Code != http.StatusOK && test.errorExpected != "" {
				var resp entity.BookResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("Should not fail: found error %v ", err)
				}

				if resp.Message != test.errorExpected {
					t.Errorf("Handler %s returned with incorrect error - got (%s) wanted (%s)", test.handler, resp.Message, test.errorExpected)
				}
			}
		})
	}

}
