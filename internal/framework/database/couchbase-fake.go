//go:build fake

package database

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/adapter/repository"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"

	"github.com/couchbase/gocb/v2"
)

const testFolderPath = "../../../tests/"

var count int = 0

// Couchbase fake
type Couchbase struct {
	Bucket  *FakeBucket
	Cluster *FakeCluster
	Force   string
}

type FakeCluster struct {
	Force string
}

type FakeBucket struct {
	Force string
}

type FakeScope struct {
	Force string
}

type FakeQuery struct {
}

type FakeCollection struct {
	Force string
}

type FakeResult struct {
	Force string
}

func NewFakeCouchbaseStorage(force string) (repository.Storage, error) {
	return &Couchbase{Bucket: &FakeBucket{Force: force}, Cluster: &FakeCluster{Force: force}}, nil
}

// NewCouchbaseStorage is here only to avoid the error on main.go
func NewCouchbaseStorage() (repository.Storage, error) {
	return nil, nil
}

// Query - inject our implementation for testing
func (fs *FakeScope) Query(_ string, _ *gocb.QueryOptions) (*FakeResult, error) {
	if fs.Force == "query-error" {
		return &FakeResult{}, errors.New("forced query error")
	}
	return &FakeResult{Force: fs.Force}, nil
}

// Scope - override the original gocb implementation
func (fb *FakeBucket) Scope(_ string) *FakeScope {
	return &FakeScope{Force: fb.Force}
}

// Collection returns an instance of a collection.
func (fs *FakeScope) Collection(_ string) *FakeCollection {
	return &FakeCollection{Force: fs.Force}
}

// Get - override the original golang implementation
func (fc *FakeCollection) Get(_ string, _ interface{}) (*FakeResult, error) {
	if fc.Force == "true" {
		return &FakeResult{}, errors.New("forced collection error")
	}
	if fc.Force == "not-found-error" {
		return &FakeResult{}, errors.New("document not found")
	}
	return &FakeResult{Force: fc.Force}, nil
}

// Next - override the original golang implementation
func (fr *FakeResult) Next() bool {
	if count < 2 {
		count++
		return true
	}
	count = 0
	return false
}

// Row - override the original golang implementation
func (fr *FakeResult) Row(_ interface{}) error {
	return nil
}

// One - override the original golang implementation
func (fr *FakeResult) One(_ interface{}) error {
	if fr.Force == "read-error" {
		return errors.New("no result was available")
	}
	return nil
}

// Content - override the original golang implementation
func (fr *FakeResult) Content(ptr interface{}) error {
	var book entity.Book
	if fr.Force == "error" {
		return errors.New("forced content error")
	}

	_, ok := ptr.(*entity.Book)
	if ok {
		data, _ := os.ReadFile(testFolderPath + "book.json")
		_ = json.Unmarshal(data, &book)
		*ptr.(*entity.Book) = book
	}

	return nil
}

// Close - do we need to explain ?
func (fr *FakeResult) Close() error {
	if fr.Force == "close-error" {
		return errors.New("forced close error")
	}
	return nil
}

// Upsert : wrapper function for couchbase upsert
func (fc *FakeCollection) Upsert(_ string, _ interface{}, _ *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if fc.Force == "error" || fc.Force == "update-error" {
		return &gocb.MutationResult{}, errors.New("forced collection upsert error")
	}
	return &gocb.MutationResult{}, nil
}
