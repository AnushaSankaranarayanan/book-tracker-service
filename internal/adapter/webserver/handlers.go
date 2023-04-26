package webserver

import (
	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var l = logrus.StandardLogger()

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

func (s *Server) ListBooks(c *gin.Context) {
	books, err := s.Services.BookTracker.ListBooks()
	if err != nil {
		l.Errorf("GetBooks error %s", err.Error())
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "failed to get books.Refer to logs for more details"))
		return
	}
	c.JSON(http.StatusOK, entity.NewBookResponse(http.StatusOK, "books retrieval successful", nil, books))
}

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

func (s *Server) UpdateBook(c *gin.Context) {
	var book entity.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		l.Errorf("UpdateBook invalid request. Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, entity.NewGenericResponse(http.StatusBadRequest, err.Error()))
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

func handleErrorTypes(c *gin.Context, err error) {
	switch err.(type) {
	case entity.NotFoundError:
		c.JSON(http.StatusNotFound, entity.NewGenericResponse(http.StatusNotFound, err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, entity.NewGenericResponse(http.StatusInternalServerError, "operation failed.Refer to logs for more details"))
	}
}
