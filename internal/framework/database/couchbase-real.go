//go:build real

package database

import (
	"github.com/couchbase/gocb/v2"
	"os"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/adapter/repository"
)

type Couchbase struct {
	Bucket  *gocb.Bucket
	Cluster *gocb.Cluster
}

func NewCouchbaseStorage() (repository.Storage, error) {
	opts := gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	}

	if os.Getenv("ENABLE_DB_VERBOSE_LOGGING") == "true" {
		gocb.SetLogger(gocb.VerboseStdioLogger())
	}
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		return nil, err
	}

	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"))

	return &Couchbase{Bucket: bucket, Cluster: cluster}, nil
}
