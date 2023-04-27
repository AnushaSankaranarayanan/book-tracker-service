package webserver

import (
	"fmt"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/consts"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	unreadStatus     = "UNREAD"
	inProgressStatus = "IN PROGRESS"
	finishedStatus   = "FINISHED"
)

var l = logrus.StandardLogger()

// AddBook - checks incoming request and add the book to DB
func (s *Server) AddBook(c *gin.Context) {
	var book entity.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		l.Errorf("CreateBook invalid request. Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, entity.NewGenericResponse(http.StatusBadRequest, err.Error()))
		return
	}

	err := s.Services.BookTracker.AddBook(book)
	if err != nil {
		l.Errorf("AddBook error %s. Request payload %+v", err.Error(), book)
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "failed to save book.Refer to logs for more details"))
		return
	}

	c.JSON(http.StatusOK, entity.NewGenericResponse(http.StatusOK, "book creation successful"))
}

// ListBooks - checks incoming parameters(sortKey) and fetches books from DB
func (s *Server) ListBooks(c *gin.Context) {
	sortKey := c.Query(consts.SortKey)

	if !sortKeyValid(sortKey) {
		msg := fmt.Sprintf("Invalid sort key. Expected: %s or %s", consts.Status, consts.Title)
		l.Errorf("GetBooks error: %s", msg)
		c.JSON(http.StatusBadRequest, entity.NewGenericResponse(http.StatusBadRequest, msg))
		return
	}

	books, err := s.Services.BookTracker.ListBooks(sortKey)
	if err != nil {
		l.Errorf("GetBooks error %s", err.Error())
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "failed to get books.Refer to logs for more details"))
		return
	}

	c.JSON(http.StatusOK, entity.NewBookResponse(http.StatusOK, "books retrieval successful", nil, books))
}

// GetBook - gets the book with id from the DB(ISBN)
func (s *Server) GetBook(c *gin.Context) {
	bookId, _ := c.Params.Get("id")
	book, err := s.Services.BookTracker.GetBook(bookId)
	if err != nil {
		l.Errorf("GetBook error %s. Request ISBN %s", err.Error(), bookId)
		handleErrorTypes(c, err)
		return
	}

	c.JSON(http.StatusOK, entity.NewBookResponse(http.StatusOK, "book retrieval successful", book, nil))
}

// UpdateBook - checks incoming request and updates the book to DB
func (s *Server) UpdateBook(c *gin.Context) {
	var book entity.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		l.Errorf("UpdateBook invalid request. Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, entity.NewGenericResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if !statusValid(book.Status) {
		msg := fmt.Sprintf("Invalid status key. Expected one of %s, %s, %s", unreadStatus, inProgressStatus, finishedStatus)
		l.Errorf("UpdateBook error: %s", msg)
		c.JSON(http.StatusBadRequest, entity.NewGenericResponse(http.StatusBadRequest, msg))
		return
	}

	err := s.Services.BookTracker.UpdateBook(book)
	if err != nil {
		l.Errorf("UpdateBook error %s. Request payload %+v", err.Error(), book)
		handleErrorTypes(c, err)
		return
	}

	c.JSON(http.StatusOK, entity.NewGenericResponse(http.StatusOK, "book updated successfully"))
}

// GroupBooksByGenre - lists the genres and books associated with each genre
func (s *Server) GroupBooksByGenre(c *gin.Context) {
	genres, err := s.Services.BookTracker.GroupBooksByGenre()
	if err != nil {
		l.Errorf("GetBooks error %s", err.Error())
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "failed to get books.Refer to logs for more details"))
		return
	}

	c.JSON(http.StatusOK, entity.NewGroupByGenreResponse(http.StatusOK, "books retrieval successful", genres))
}

func sortKeyValid(sortKey string) bool {
	return sortKey == "" || strings.ToLower(sortKey) == consts.Title || strings.ToLower(sortKey) == consts.Status
}

func statusValid(status string) bool {
	return status == "" ||
		strings.ToUpper(status) == unreadStatus ||
		strings.ToUpper(status) == inProgressStatus ||
		strings.ToUpper(status) == finishedStatus
}

func handleErrorTypes(c *gin.Context, err error) {
	switch err.(type) {
	case entity.NotFoundError:
		c.JSON(http.StatusNotFound, entity.NewGenericResponse(http.StatusNotFound, err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "operation failed.Refer to logs for more details"))
	}
}
