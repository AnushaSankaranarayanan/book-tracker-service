package repository

import "github.com/anushasankaranarayanan/book-tracker-service/internal/entity"

type Storage interface {
	Upsert(string, interface{}) error
	GetAll() ([]entity.Book, error)
	Get(string) (*entity.Book, error)
}
