package service

import (
	"fmt"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/consts"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
	"time"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
)

const (
	documentNotFoundError = "document not found"
)

var l = logrus.StandardLogger()

type BookTracker interface {
	AddBook(entity.Book) error
	UpdateBook(entity.Book) error
	ListBooks(string) ([]entity.Book, error)
	GetBook(string) (*entity.Book, error)
}

type BookRepository interface {
	Upsert(string, interface{}) error
	GetAll() ([]entity.Book, error)
	Get(string) (*entity.Book, error)
	GetByScope(string) (entity.Book, error)
}

type bookTracker struct {
	storage BookRepository
}

func NewBookTracker(tr BookRepository) BookTracker {
	return &bookTracker{storage: tr}
}

func (svc *bookTracker) AddBook(book entity.Book) error {
	book.SetTrackingDetails()
	err := svc.storage.Upsert(book.ISBN, book)

	if err != nil {
		return err
	}

	l.Infof("book %s inserted into couchbase successfully", book.Title)

	return nil
}

func (svc *bookTracker) UpdateBook(book entity.Book) error {
	id := book.ISBN
	book.Updated = time.Now().Unix()

	_, err := svc.GetBook(id)
	if err != nil {
		return err
	}

	err = svc.storage.Upsert(id, book)
	if err != nil {
		return err
	}

	l.Infof("book %s updated successfully", id)
	return nil
}

func (svc *bookTracker) ListBooks(sortKey string) ([]entity.Book, error) {
	books, err := svc.storage.GetAll()
	if err != nil {
		return nil, err
	}
	sortResponse(sortKey, books)
	return books, nil
}

func (svc *bookTracker) GetBook(id string) (*entity.Book, error) {
	book, err := svc.storage.Get(id)
	if err != nil && strings.Contains(err.Error(), documentNotFoundError) {
		return nil, entity.NotFoundError{Message: fmt.Sprintf("book with id %s not found", id)}
	}
	return book, err
}

func (svc *bookTracker) GetByScope(scope string) (entity.Book, error) {
	return svc.storage.GetByScope(scope)
}

func sortResponse(sortKey string, books []entity.Book) {
	if strings.ToLower(sortKey) == consts.Title {
		sort.Slice(books, func(i, j int) bool {
			return books[i].Title < books[j].Title
		})
	}
	if strings.ToLower(sortKey) == consts.Status {
		sort.Slice(books, func(i, j int) bool {
			return books[i].Status < books[j].Status
		})
	}
}
