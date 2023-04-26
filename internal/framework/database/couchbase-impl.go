package database

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/sirupsen/logrus"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
)

const (
	defaultScope   = "_default"
	bookCollection = "book"
)

var l = logrus.StandardLogger()

// Get - wrapper to get a book resource
func (c *Couchbase) Get(id string) (*entity.Book, error) {
	var book entity.Book

	collection := c.Bucket.Scope(defaultScope).Collection(bookCollection)
	result, err := collection.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("get error:%s", err.Error())
	}
	err = result.Content(&book)
	if err != nil {
		return nil, fmt.Errorf("get content error:%s", err.Error())
	}
	return &book, nil
}

// GetAll - wrapper to list all book resources
func (c *Couchbase) GetAll() ([]entity.Book, error) {
	var books []entity.Book

	query := "select raw b from book b"

	l.Tracef("Function GetAll %s", query)
	res, err := c.Bucket.Scope(defaultScope).Query(query, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAll query error:%s", err.Error())
	}

	l.Tracef("Function GetAll next()")
	for res.Next() {
		var row entity.Book
		err = res.Row(&row)
		l.Tracef("Function GetAll data %+v", row)
		books = append(books, row)
	}
	if err = res.Close(); err != nil {
		return nil, fmt.Errorf("GetAll result close error:%s", err.Error())
	}

	return books, nil
}

// Upsert :  wrapper to update a book resource
func (c *Couchbase) Upsert(key string, value interface{}) error {
	opts := &gocb.UpsertOptions{}
	collection := c.Bucket.Scope(defaultScope).Collection(bookCollection)
	_, err := collection.Upsert(key, value, opts)
	if err != nil {
		return fmt.Errorf("Upsert error:%s", err.Error())
	}
	return nil
}

// GetByScope - wrapper to get a book by scope
func (c *Couchbase) GetByScope(scope string) (entity.Book, error) {
	var book entity.Book

	query := "select raw t from book t where t.`scope` = $1"

	l.Tracef("Function GetByScope %s", query)
	res, err := c.Bucket.Scope(defaultScope).Query(query, &gocb.QueryOptions{PositionalParameters: []interface{}{scope}})
	if err != nil {
		return book, fmt.Errorf("GetByScope query error:%s", err.Error())
	}

	err = res.One(&book)
	if err != nil {
		return book, fmt.Errorf("GetByScope result extract error:%s", err.Error())
	}

	if err = res.Close(); err != nil {
		return book, fmt.Errorf("GetByScope result close error:%s", err.Error())
	}

	l.Tracef("Book retrieved from DB %+v", book)
	return book, nil
}
