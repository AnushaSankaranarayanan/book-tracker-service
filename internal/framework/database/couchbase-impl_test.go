//go:build fake

package database

import (
	"errors"
	"testing"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/entity"
)

const (
	upsertMethod            = "Upsert"
	getMethod               = "Get"
	getAllMethod            = "ListBooks"
	getByScopeMethod        = "GetByScope"
	newFakeCouchbaseStorage = "NewFakeCouchbaseStorage"
	newCouchbaseStorage     = "NewCouchbaseStorage"
)

func TestCouchbaseImpl(t *testing.T) {

	tests := []struct {
		testName  string
		errorFlag string
		arg       string
		method    string
		expected  error
	}{

		{
			"Upsert: should pass",
			"",
			"ISBN-01",
			upsertMethod,
			nil,
		},
		{
			"Upsert: should fail (force update-error)",
			"update-error",
			"ISBN-01",
			upsertMethod,
			errors.New("Upsert error:forced collection upsert error"),
		},
		{
			"GetById: should pass",
			"",
			"ISBN-01",
			getMethod,
			nil,
		},
		{
			"GetById: should fail (collection error)",
			"true",
			"ISBN-01",
			getMethod,
			errors.New("get error:forced collection error"),
		},
		{
			"GetById: should fail (content error)",
			"error",
			"ISBN-01",
			getMethod,
			errors.New("get content error:forced content error"),
		},
		{
			"GetAll: should pass",
			"",
			"",
			getAllMethod,
			nil,
		},
		{
			"GetAll query error:forced query error",
			"query-error",
			"",
			getAllMethod,
			errors.New("GetAll query error:forced query error"),
		},
		{
			"GetAll result close error:forced close error",
			"close-error",
			"",
			getAllMethod,
			errors.New("GetAll result close error:forced close error"),
		},
		{
			"GetByScope: should pass",
			"",
			"ISBN-01",
			getByScopeMethod,
			nil,
		},
		{
			"GetByScope: should fail (query-error)",
			"query-error",
			"ISBN-01",
			getByScopeMethod,
			errors.New("GetByScope query error:forced query error"),
		},
		{
			"GetByScope: should fail (read-error)",
			"read-error",
			"ISBN-01",
			getByScopeMethod,
			errors.New("GetByScope result extract error:no result was available"),
		}, {
			"GetByScope: should fail (close-error)",
			"close-error",
			"ISBN-01",
			getByScopeMethod,
			errors.New("GetByScope result close error:forced close error"),
		},
		{
			"NewFakeCouchbaseStorage: should pass",
			"",
			"",
			newFakeCouchbaseStorage,
			nil,
		},
		{
			"NewCouchbaseSorage: should pass",
			"",
			"",
			newCouchbaseStorage,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockCouchbase := &Couchbase{Bucket: &FakeBucket{Force: test.errorFlag}, Cluster: &FakeCluster{}}

			var err error
			switch test.method {
			case upsertMethod:
				err = mockCouchbase.Upsert(test.arg, entity.Book{})
			case getMethod:
				_, err = mockCouchbase.Get(test.arg)
			case getAllMethod:
				_, err = mockCouchbase.GetAll()
			case getByScopeMethod:
				_, err = mockCouchbase.GetByScope(test.arg)
			case newFakeCouchbaseStorage:
				_, err = NewFakeCouchbaseStorage("")
			case newCouchbaseStorage:
				_, err = NewCouchbaseStorage()
			}

			if err == nil && test.expected != err {
				t.Errorf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", test.method, err, nil)
			}

			if test.expected != nil && test.expected.Error() != err.Error() {
				t.Errorf("Function (%s) assert (error type is different from expected) -  got (%s) wanted (%s)", test.method, err.Error(), test.expected.Error())
			}
		})
	}
}
