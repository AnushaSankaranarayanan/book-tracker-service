package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/adapter/webserver"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/framework/database"
	"github.com/anushasankaranarayanan/book-tracker-service/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error %s\\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	logger := logrus.StandardLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Info(".env file not detected.... falling through to Kubernetes ✿✿")
	}

	cbStorage, err := database.NewCouchbaseStorage()
	if err != nil {
		logger.Errorf("Couchbase connection error: %v", err)
		return err
	}

	bookTrackingSvc := service.NewBookTracker(cbStorage)

	services := webserver.Services{
		BookTracker: bookTrackingSvc,
	}

	server := webserver.NewServer(services)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}
